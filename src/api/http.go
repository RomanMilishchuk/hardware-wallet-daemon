package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/cors"
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120

	// ContentTypeJSON json content type header
	ContentTypeJSON = "application/json"
	// ContentTypeForm form data content type header
	ContentTypeForm = "application/x-www-form-urlencoded"

	apiVersion1 = "v1"
)

var (
	logger = logging.MustGetLogger("daemon-api")
)

type muxConfig struct {
	host               string
	disableHeaderCheck bool
	hostWhitelist      []string
}

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	DisableHeaderCheck bool
	HostWhitelist      []string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// HTTPResponse represents the http response struct
type HTTPResponse struct {
	Error *HTTPError  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// ReceivedHTTPResponse parsed is a Parsed HTTPResponse
type ReceivedHTTPResponse struct {
	Error *HTTPError      `json:"error,omitempty"`
	Data  json.RawMessage `json:"data"`
}

// HTTPError is included in an HTTPResponse
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewHTTPErrorResponse returns an HTTPResponse with the Error field populated
func NewHTTPErrorResponse(code int, msg string) HTTPResponse {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return HTTPResponse{
		Error: &HTTPError{
			Code:    code,
			Message: msg,
		},
	}
}

func writeHTTPResponse(w http.ResponseWriter, resp HTTPResponse) {
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		wh.Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", ContentTypeJSON)

	if resp.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if resp.Error.Code < 400 || resp.Error.Code >= 600 {
			logger.Critical().Errorf("writeHTTPResponse invalid error status code: %d", resp.Error.Code)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(resp.Error.Code)
		}
	}

	if _, err := w.Write(out); err != nil {
		logger.WithError(err).Error("http Write failed")
	}
}

// Serve serves the web interface on the configured host
func (s *Server) Serve() error {
	defer close(s.done)

	if err := s.server.Serve(s.listener); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

// Shutdown closes the HTTP service. This can only be called after Serve or ServeHTTPS has been called.
func (s *Server) Shutdown() {
	if s == nil {
		return
	}

	logger.Info("Shutting down web interface")
	defer logger.Info("Web interface shut down")
	if err := s.listener.Close(); err != nil {
		logger.WithError(err).Warning("s.listener.Close() error")
	}
	<-s.done
}

func create(host string, c Config, gateway *Gateway) (*Server, error) {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaultIdleTimeout
	}

	mc := muxConfig{
		host: host,
		disableHeaderCheck: c.DisableHeaderCheck,
		hostWhitelist: c.HostWhitelist,
	}

	srvMux := newServerMux(mc, gateway.USBDevice, gateway.EmulatorDevice)

	srv := &http.Server{
		Handler:      srvMux,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	return &Server{
		server: srv,
		done:   make(chan struct{}),
	}, nil
}

func Create(host string, c Config, gateway *Gateway) (*Server, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s, err := create(host, c, gateway)
	if err != nil {
		if closeErr := s.listener.Close(); closeErr != nil {
			logger.WithError(err).Warning("s.listener.Close() error")
		}
		return nil, err
	}

	s.listener = listener

	return s, nil
}

func newServerMux(c muxConfig, usbGateway, emulatorGateway Gatewayer) *http.ServeMux {
	mux := http.NewServeMux()

	allowedOrigins := []string{fmt.Sprintf("http://%s", c.host)}
	for _, s := range c.hostWhitelist {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s", s))
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:     allowedOrigins,
		Debug:              false,
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		AllowCredentials:   false, // credentials are not used, but it would be safe to enable if necessary
		OptionsPassthrough: false,
	})

	headerCheck := func(apiVersion, host string, hostWhitelist []string, handler http.Handler) http.Handler {
		handler = originRefererCheck(apiVersion, host, hostWhitelist, handler)
		handler = hostCheck(apiVersion, host, hostWhitelist, handler)
		return handler
	}

	webHandlerWithOptionals := func(apiVersion, endpoint string, handlerFunc http.Handler, checkHeaders bool) {
		handler := wh.ElapsedHandler(logger, handlerFunc)

		handler = corsHandler.Handler(handler)

		if checkHeaders {
			handler = headerCheck(apiVersion, c.host, c.hostWhitelist, handler)
		}

		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandler := func(apiVersion, endpoint string, handler http.Handler) {
		handler = wh.ElapsedHandler(logger, handler)

		// mux.Handle("/api"+endpoint, handler)

		webHandlerWithOptionals(apiVersion1, endpoint, handler, !c.disableCSP)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandler(apiVersion1, "/api/v1"+endpoint, handler)
	}

	// hw wallet endpoints
	webHandlerV1("/generateAddresses", generateAddresses(usbGateway))
	webHandlerV1("/applySettings", applySettings(usbGateway))
	webHandlerV1("/backup", backup(usbGateway))
	webHandlerV1("/cancel", cancel(usbGateway))
	webHandlerV1("/checkMessageSignature", checkMessageSignature(usbGateway))
	webHandlerV1("/features", features(usbGateway))
	webHandlerV1("/generateMnemonic", generateMnemonic(usbGateway))
	webHandlerV1("/recovery", recovery(usbGateway))
	webHandlerV1("/setMnemonic", setMnemonic(usbGateway))
	webHandlerV1("/setPinCode", setPinCode(usbGateway))
	webHandlerV1("/signMessage", signMessage(usbGateway))
	webHandlerV1("/transactionSign", transactionSign(usbGateway))
	webHandlerV1("/wipe", wipe(usbGateway))
	webHandlerV1("/intermediate/pinmatrix", PinMatrixRequestHandler(usbGateway))
	webHandlerV1("/intermediate/passphrase", PassphraseRequestHandler(usbGateway))
	webHandlerV1("/intermediate/word", WordRequestHandler(usbGateway))

	// emulator endpoints
	webHandlerV1("/emulator/generateAddresses", generateAddresses(emulatorGateway))
	webHandlerV1("/emulator/applySettings", applySettings(emulatorGateway))
	webHandlerV1("/emulator/backup", backup(emulatorGateway))
	webHandlerV1("/emulator/cancel", cancel(emulatorGateway))
	webHandlerV1("/emulator/checkMessageSignature", checkMessageSignature(emulatorGateway))
	webHandlerV1("/emulator/features", features(emulatorGateway))
	webHandlerV1("/emulator/generateMnemonic", generateMnemonic(emulatorGateway))
	webHandlerV1("/emulator/recovery", recovery(emulatorGateway))
	webHandlerV1("/emulator/setMnemonic", setMnemonic(emulatorGateway))
	webHandlerV1("/emulator/setPinCode", setPinCode(emulatorGateway))
	webHandlerV1("/emulator/signMessage", signMessage(emulatorGateway))
	webHandlerV1("/emulator/transactionSign", transactionSign(emulatorGateway))
	webHandlerV1("/emulator/wipe", wipe(emulatorGateway))
	webHandlerV1("/emulator/intermediate/pinmatrix", PinMatrixRequestHandler(emulatorGateway))
	webHandlerV1("/emulator/intermediate/passphrase", PassphraseRequestHandler(emulatorGateway))
	webHandlerV1("/emulator/intermediate/word", WordRequestHandler(emulatorGateway))

	return mux
}

func parseBoolFlag(v string) (bool, error) {
	if v == "" {
		return false, nil
	}

	return strconv.ParseBool(v)
}

type IntermediateResponse struct {
	RequestType string `json:"request_type"`
}

func HandleFirmwareResponseMessages(w http.ResponseWriter, r *http.Request, gateway Gatewayer, msg wire.Message) {
	switch msg.Kind {
	case uint16(messages.MessageType_MessageType_PinMatrixRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "PinMatrixRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_PassphraseRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "PassPhraseRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_WordRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "WordRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_ButtonRequest):
		msg, err := gateway.ButtonAck()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnauthorized, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	case uint16(messages.MessageType_MessageType_Failure):
		failureMsg, err := deviceWallet.DecodeFailMsg(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		resp := NewHTTPErrorResponse(http.StatusConflict, failureMsg)
		writeHTTPResponse(w, resp)
		return
	case uint16(messages.MessageType_MessageType_Success):
		successMsg, err := deviceWallet.DecodeSuccessMsg(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnauthorized, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: successMsg,
		})
	// AddressGen Response
	case uint16(messages.MessageType_MessageType_ResponseSkycoinAddress):
		addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: GenerateAddressesResponse{
				Addresses: addresses,
			},
		})
	// Features Response
	case uint16(messages.MessageType_MessageType_Features):
		features := &messages.Features{}
		err := proto.Unmarshal(msg.Data, features)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: FeaturesResponse{
				Features: features,
			},
		})
	// SignMessage Response
	case uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage):
		signature, err := deviceWallet.DecodeResponseSkycoinSignMessage(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: SignMessageResponse{
				Signature: signature,
			},
		})
	// TransactionSign Response
	case uint16(messages.MessageType_MessageType_ResponseTransactionSign):
		signatures, err := deviceWallet.DecodeResponseTransactionSign(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: TransactionSignResponse{
				Signatures: signatures,
			},
		})
	default:
		resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("recevied unexpected response message type: %s", messages.MessageType(msg.Kind)))
		writeHTTPResponse(w, resp)
	}
}

type PinMatrixRequest struct {
	Pin string `json:"pin"`
}

func PinMatrixRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req PinMatrixRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		msg, err := gateway.PinMatrixAck(req.Pin)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

type PassPhraseRequest struct {
	Passphrase string `json:"passphrase"`
}

func PassphraseRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req PassPhraseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		msg, err := gateway.PassphraseAck(req.Passphrase)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

type WordRequest struct {
	Word string `json:"word"`
}

func WordRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req WordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		msg, err := gateway.WordAck(req.Word)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

func newStrPtr(s string) *string {
	return &s
}

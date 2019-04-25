package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gogo/protobuf/proto"
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	wh "github.com/skycoin/skycoin/src/util/http"
)

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

// HandleFirmwareResponseMessages handles response messages from the firmware
func HandleFirmwareResponseMessages(w http.ResponseWriter, gateway Gatewayer, msg wire.Message) {
	for {
		switch msg.Kind {
		case uint16(messages.MessageType_MessageType_PinMatrixRequest):
			writeHTTPResponse(w, HTTPResponse{
				Data: "PinMatrixRequest",
			})
			return
		case uint16(messages.MessageType_MessageType_PassphraseRequest):
			writeHTTPResponse(w, HTTPResponse{
				Data: "PassPhraseRequest",
			})
			return
		case uint16(messages.MessageType_MessageType_WordRequest):
			writeHTTPResponse(w, HTTPResponse{
				Data: "WordRequest",
			})
			return
		case uint16(messages.MessageType_MessageType_ButtonRequest):
			var err error
			msg, err = gateway.ButtonAck()
			if err != nil {
				logger.Error(err.Error())
				resp := NewHTTPErrorResponse(http.StatusUnauthorized, err.Error())
				writeHTTPResponse(w, resp)
			}
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
			return

		// AddressGen Response
		case uint16(messages.MessageType_MessageType_ResponseSkycoinAddress):
			addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}

			writeHTTPResponse(w, HTTPResponse{
				Data: addresses,
			})
			return

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
				Data: features,
			})
			return

		// SignMessage Response
		case uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage):
			signature, err := deviceWallet.DecodeResponseSkycoinSignMessage(msg)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}

			writeHTTPResponse(w, HTTPResponse{
				Data: signature,
			})
			return

		// TransactionSign Response
		case uint16(messages.MessageType_MessageType_ResponseTransactionSign):
			signatures, err := deviceWallet.DecodeResponseTransactionSign(msg)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}

			writeHTTPResponse(w, HTTPResponse{
				Data: &signatures,
			})
			return
		default:
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("recevied unexpected response message type: %s", messages.MessageType(msg.Kind)))
			writeHTTPResponse(w, resp)
			return
		}
	}
}

func serialize(gateway Gatewayer) (func(), error) {
	select {
	// acquire the lock
	case ongoingOperation <- struct{}{}:
		// defer func to release the lock
		return func() {
			<-ongoingOperation
		}, nil
	default:
		// cancel any on-going operation if a lock is acquired already
		_, err := gateway.Cancel()
		if err != nil {
			return nil, err
		}

		ongoingOperation <- struct{}{}
		return func() {
			<-ongoingOperation
		}, nil
	}
}

func parseBoolFlag(v string) (bool, error) {
	if v == "" {
		return false, nil
	}

	return strconv.ParseBool(v)
}

func newStrPtr(s string) *string {
	return &s
}

func newUint32Ptr(n uint32) *uint32 {
	return &n
}

func newBoolPtr(b bool) *bool {
	return &b
}

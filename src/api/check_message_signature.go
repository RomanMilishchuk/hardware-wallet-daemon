package api

import (
	"encoding/json"
	"net/http"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	"github.com/skycoin/skycoin/src/cipher"
)

// CheckMessageSignatureRequest is request data for /api/v1/check_message_signature
type CheckMessageSignatureRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

// URI: /api/v1/checkMessageSignature
// Method: POST
// Content-Type: application/json
// Args: JSON Body
func checkMessageSignature(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// allow only one request at a time
		closeFunc, err := serialize(gateway)
		if err != nil {
			logger.Error("serialize failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer closeFunc()

		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req CheckMessageSignatureRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if req.Address == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "address is required")
			writeHTTPResponse(w, resp)
			return
		}

		_, err = cipher.DecodeBase58Address(req.Address)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if req.Signature == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "signature is required")
			writeHTTPResponse(w, resp)
			return
		}

		if req.Message == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "message is required")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("checkMessageSignature failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.CheckMessageSignature(req.Message, req.Signature, req.Address)
		if err != nil {
			logger.Errorf("checkMessageSignature failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, gateway, msg)
	}
}

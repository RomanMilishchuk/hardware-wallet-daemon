package api

import (
	"encoding/json"
	"net/http"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

// SignMessageRequest is request data for /api/v1/sign_message
type SignMessageRequest struct {
	AddressN int    `json:"address_n"`
	Message  string `json:"message"`
}

// SignMessageResponse is data returned by POST /api/v1/sign_message
type SignMessageResponse struct {
	Signature string `json:"signature"`
}

// URI: /api/v1/signMessage
// Method: POST
// Args: JSON Body
func signMessage(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var req SignMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if req.AddressN < 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be negative")
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
			err := gateway.SetAutoPressButton(true, deviceWallet.ButtonRight)
			if err != nil {
				logger.Error("generateAddress failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.SignMessage(req.AddressN, req.Message)
		if err != nil {
			logger.Errorf("signMessage failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

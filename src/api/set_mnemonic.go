package api

import (
	"encoding/json"
	"net/http"
)

// SetMnemonicRequest is request data for /api/v1/set_mnemonic
type SetMnemonicRequest struct {
	Mnemonic string `json:"mnemonic"`
}

// URI: /api/v1/set_mnemonic
// Method: POST
// Args: JSON Body
func setMnemonic(gateway Gatewayer) http.HandlerFunc {
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

		var req SetMnemonicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		// TODO(therealssj): add mnemonic check?

		msg, err := gateway.SetMnemonic(req.Mnemonic)
		if err != nil {
			logger.Errorf("setMnemonic failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

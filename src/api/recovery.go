package api

import (
	"encoding/json"
	"net/http"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

// RecoveryRequest is request data for /api/v1/recovery
type RecoveryRequest struct {
	WordCount     uint32 `json:"word_count"`
	UsePassphrase bool   `json:"use_passphrase"`
	DryRun        bool   `json:"dry_run"`
}

// URI: /api/v1/recovery
// Method: POST
// Args: JSON Body
func recovery(gateway Gatewayer) http.HandlerFunc {
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

		var req RecoveryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

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

		msg, err := gateway.Recovery(req.WordCount, req.UsePassphrase, req.DryRun)
		if err != nil {
			logger.Errorf("recovery failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, gateway, msg)
	}
}

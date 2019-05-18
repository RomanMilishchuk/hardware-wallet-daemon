package api

import (
	"encoding/json"
	"net/http"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

// GenerateMnemonicRequest is request data for /api/v1/generate_mnemonic
type GenerateMnemonicRequest struct {
	WordCount     uint32 `json:"word_count"`
	UsePassphrase bool   `json:"use_passphrase"`
}

// URI: /api/v1/generate_mnemonic
// Method: POST
// Args: JSON Body
func generateMnemonic(gateway Gatewayer) http.HandlerFunc {
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

		var req GenerateMnemonicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if req.WordCount != 12 && req.WordCount != 24 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "word count must be 12 or 24")
				writeHTTPResponse(w, resp)
				return
			}
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("generateMnemonic failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.GenerateMnemonic(req.WordCount, req.UsePassphrase)
		if err != nil {
			logger.Errorf("generateMnemonic failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, gateway, msg)
	}
}

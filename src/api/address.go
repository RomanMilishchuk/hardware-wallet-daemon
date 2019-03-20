package api

import (
	"encoding/json"
	"net/http"
)

// GenerateAddressesRequest is request data for /api/v1/generateAddresses
type GenerateAddressesRequest struct {
	AddressN       int  `json:"address_n"`
	StartIndex     int  `json:"start_index"`
	ConfirmAddress bool `json:"confirm_address"`
}

// GenerateAddressesResponse is returned by POST /api/v1/generateAddresses
type GenerateAddressesResponse struct {
	Addresses []string `json:"addresses"`
}

// generateAddresses gene
// URI: /api/v1/generateAddresses
// Method: POST
// Args: JSON Body
func generateAddresses(gateway Gatewayer) http.HandlerFunc {
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

		var req GenerateAddressesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if req.AddressN == 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be 0")
			writeHTTPResponse(w, resp)
			return
		}

		if req.AddressN < 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be negative")
			writeHTTPResponse(w, resp)
			return
		}

		if req.StartIndex < 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "start_index cannot be negative")
			writeHTTPResponse(w, resp)
			return
		}
		
		msg, err := gateway.AddressGen(req.AddressN, req.StartIndex, req.ConfirmAddress)
		if err != nil {
			logger.Error("generateAddress failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

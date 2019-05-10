package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	"github.com/stretchr/testify/require"
)

// TODO(therealssj): add more tests

func TestFirmwareUpdate(t *testing.T) {
	postData :=
		`--xxx
Content-Disposition: form-data; name="file"

value1
--xxx
Content-Disposition: form-data; name="file"

value2
--xxx
Content-Disposition: form-data; name="file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xx--
`

	cases := []struct {
		name         string
		method       string
		status       int
		data         string
		header       string
		emulator     bool
		httpResponse HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "404 - Not Found",
			method:       http.MethodGet,
			status:       http.StatusNotFound,
			httpResponse: NewHTTPErrorResponse(http.StatusNotFound, ""),
			emulator:     true,
		},

		{
			name:         "400 - EOF",
			method:       http.MethodPut,
			status:       http.StatusBadRequest,
			header:       `multipart/form-data; boundary=xxx`,
			data:         postData,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "unexpected EOF"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/firmware_update"
			gateway := &MockGatewayer{}

			req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, nil)
			require.NoError(t, err)

			if tc.data != "" {
				req.Header.Set("Content-Type", tc.header)
				req.Body = ioutil.NopCloser(strings.NewReader(tc.data))
			}

			rr := httptest.NewRecorder()

			config := defaultMuxConfig()
			if tc.emulator {
				config.mode = skyWallet.DeviceTypeEmulator
			}

			handler := newServerMux(config, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if !tc.emulator {
				var rsp ReceivedHTTPResponse
				err = json.NewDecoder(rr.Body).Decode(&rsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Error, rsp.Error)

				if rsp.Data == nil {
					require.Nil(t, tc.httpResponse.Data)
				} else {
					require.NotNil(t, tc.httpResponse.Data)
				}
			} else {
				// check that it returns 404
				require.Equal(t, tc.httpResponse.Error.Code, rr.Code)
			}
		})
	}
}

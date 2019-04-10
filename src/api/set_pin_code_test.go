package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	"github.com/stretchr/testify/require"
)

func TestSetPinCode(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("setpincode success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                    string
		method                  string
		status                  int
		gatewaySetPinCodeResult wire.Message
		httpResponse            HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodPost,
			status: http.StatusConflict,
			gatewaySetPinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			gatewaySetPinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "setpincode success msg",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/set_pin_code"
			gateway := &MockGatewayer{}

			gateway.On("ChangePin").Return(tc.gatewaySetPinCodeResult, nil)

			req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var resp string
				err = json.Unmarshal(rsp.Data, &resp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(string), resp)
			}
		})
	}
}

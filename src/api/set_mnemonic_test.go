package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/stretchr/testify/require"
)

func TestSetMnemonic(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("setmnemonic success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                     string
		method                   string
		status                   int
		contentType              string
		httpBody                 string
		gatewaySetMnemonicResult wire.Message
		httpResponse             HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  ContentTypeForm,
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:   "422 - invalid bip39 seed",
			method: http.MethodPost,
			status: http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &SetMnemonicRequest{
				Mnemonic: "foo bar foo bar",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "seed is not a valid bip39 seed"),
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodPost,
			status: http.StatusConflict,
			gatewaySetMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpBody: toJSON(t, &SetMnemonicRequest{
				Mnemonic: "cloud flower upset remain green metal below cup stem infant art thank",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			gatewaySetMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpBody: toJSON(t, &SetMnemonicRequest{
				Mnemonic: "cloud flower upset remain green metal below cup stem infant art thank",
			}),
			httpResponse: HTTPResponse{
				Data: "setmnemonic success msg",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/set_mnemonic"
			gateway := &MockGatewayer{}

			var body SetMnemonicRequest
			err := json.Unmarshal([]byte(tc.httpBody), &body)
			if err == nil {
				gateway.On("SetMnemonic", body.Mnemonic).Return(tc.gatewaySetMnemonicResult, nil)
			}

			req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Set("Content-Type", contentType)

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

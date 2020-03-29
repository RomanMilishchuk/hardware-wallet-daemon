package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	skyWallet "github.com/SkycoinProject/hardware-wallet-go/src/skywallet"
	messages "github.com/SkycoinProject/hardware-wallet-protob/go"
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/gogo/protobuf/proto"
)

// TransactionSignRequest is request data for /api/v1/transaction_sign
type TransactionSignRequest struct {
	TransactionInputs  []TransactionInput  `json:"transaction_inputs"`
	TransactionOutputs []TransactionOutput `json:"transaction_outputs"`
}

// TransactionInput is a skycoin transaction input
type TransactionInput struct {
	Index *uint32 `json:"index"` // pointer to differentiate between 0 and nil
	Hash  string  `json:"hash"`
}

// TransactionOutput is a skycoin transaction output
type TransactionOutput struct {
	AddressIndex *uint32 `json:"address_index"` // pointer to differentiate between 0 and nil
	Address      string  `json:"address"`
	Coins        string  `json:"coins"`
	Hours        string  `json:"hours"`
}

// TransactionSignResponse is data returned by POST /api/v1/transaction_sign
type TransactionSignResponse struct {
	Signatures *[]string `json:"signatures"`
}

// URI: /api/v1/transactionSign
// Method: POST
// Args: JSON Body
func transactionSign(gateway Gatewayer) http.HandlerFunc {
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

		var req TransactionSignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if err := req.validate(); err != nil {
			logger.WithError(err).Error("invalid sign transaction request")
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		txnInputs, txnOutputs, err := req.TransactionParams()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("transactionSign failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		var signatures []string
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			signer := skyWallet.SkycoinTransactionSigner{
				Inputs:   txnInputs,
				Outputs:  txnOutputs,
				Version:  1,
				LockTime: 0,
			}
			signatures, err = gateway.GeneralTransactionSign(&signer)
			if err != nil {
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			writeHTTPResponse(w, HTTPResponse{
				Data: &signatures,
			})
		case <-errCH:
			logger.Errorf("transactionSign failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		case <-ctx.Done():
			disConnErr := gateway.Disconnect()
			if disConnErr != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				resp := NewHTTPErrorResponse(499, "Client Closed Request")
				writeHTTPResponse(w, resp)
			}
		}
	}
}

func (r *TransactionSignRequest) validate() error {
	if len(r.TransactionInputs) == 0 {
		return errors.New("inputs are required")
	}

	for _, input := range r.TransactionInputs {
		if input.Hash == "" {
			return errors.New("input hash cannot be empty")
		}
	}

	for _, output := range r.TransactionOutputs {
		if output.Address == "" {
			return errors.New("address cannot be empty")
		}

		if output.Coins == "" {
			return errors.New("coins cannot be empty")
		}

		if output.Hours == "" {
			return errors.New("hours cannot be empty")
		}
	}

	return nil
}

// TransactionParams returns params for a transaction from the request data
func (r *TransactionSignRequest) TransactionParams() ([]*messages.TxAck_TransactionType_TxInputType, []*messages.TxAck_TransactionType_TxOutputType, error) {
	var transactionInputs []*messages.TxAck_TransactionType_TxInputType
	var transactionOutputs []*messages.TxAck_TransactionType_TxOutputType

	for _, input := range r.TransactionInputs {
		var transactionInput messages.TxAck_TransactionType_TxInputType

		transactionInput.HashIn = proto.String(input.Hash)

		if input.Index != nil {
			transactionInput.AddressN = []uint32{*proto.Uint32(*input.Index)}
		} else {
			transactionInput.AddressN = nil
		}
		transactionInputs = append(transactionInputs, &transactionInput)
	}

	for _, output := range r.TransactionOutputs {
		var transactionOutput messages.TxAck_TransactionType_TxOutputType

		_, err := cipher.DecodeBase58Address(output.Address)
		if err != nil {
			return nil, nil, err
		}

		coins, err := droplet.FromString(output.Coins)
		if err != nil {
			return nil, nil, err
		}

		hours, err := strconv.ParseUint(output.Hours, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		transactionOutput.Address = proto.String(output.Address)
		transactionOutput.Coins = proto.Uint64(coins)
		transactionOutput.Hours = proto.Uint64(hours)

		if output.AddressIndex != nil {
			transactionOutput.AddressN = []uint32{*proto.Uint32(*output.AddressIndex)}
		}

		transactionOutputs = append(transactionOutputs, &transactionOutput)
	}

	return transactionInputs, transactionOutputs, nil
}

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TransactionSignRequest transaction sign request
//
// swagger:model TransactionSignRequest
type TransactionSignRequest struct {

	// transaction inputs
	// Required: true
	TransactionInputs []*TransactionInput `json:"transaction_inputs"`

	// transaction outputs
	// Required: true
	TransactionOutputs []*TransactionOutput `json:"transaction_outputs"`
}

// Validate validates this transaction sign request
func (m *TransactionSignRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTransactionInputs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransactionOutputs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionSignRequest) validateTransactionInputs(formats strfmt.Registry) error {

	if err := validate.Required("transaction_inputs", "body", m.TransactionInputs); err != nil {
		return err
	}

	for i := 0; i < len(m.TransactionInputs); i++ {
		if swag.IsZero(m.TransactionInputs[i]) { // not required
			continue
		}

		if m.TransactionInputs[i] != nil {
			if err := m.TransactionInputs[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("transaction_inputs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *TransactionSignRequest) validateTransactionOutputs(formats strfmt.Registry) error {

	if err := validate.Required("transaction_outputs", "body", m.TransactionOutputs); err != nil {
		return err
	}

	for i := 0; i < len(m.TransactionOutputs); i++ {
		if swag.IsZero(m.TransactionOutputs[i]) { // not required
			continue
		}

		if m.TransactionOutputs[i] != nil {
			if err := m.TransactionOutputs[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("transaction_outputs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *TransactionSignRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TransactionSignRequest) UnmarshalBinary(b []byte) error {
	var res TransactionSignRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

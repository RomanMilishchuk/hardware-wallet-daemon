// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SetMnemonicRequest set mnemonic request
// swagger:model SetMnemonicRequest
type SetMnemonicRequest struct {

	// mnemonic
	// Required: true
	Mnemonic *string `json:"mnemonic"`
}

// Validate validates this set mnemonic request
func (m *SetMnemonicRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMnemonic(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SetMnemonicRequest) validateMnemonic(formats strfmt.Registry) error {

	if err := validate.Required("mnemonic", "body", m.Mnemonic); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SetMnemonicRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SetMnemonicRequest) UnmarshalBinary(b []byte) error {
	var res SetMnemonicRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

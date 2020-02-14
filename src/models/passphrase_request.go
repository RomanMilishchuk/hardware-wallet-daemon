// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PassphraseRequest passphrase request
//
// swagger:model PassphraseRequest
type PassphraseRequest struct {

	// passphrase
	// Required: true
	Passphrase *string `json:"passphrase"`
}

// Validate validates this passphrase request
func (m *PassphraseRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassphrase(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PassphraseRequest) validatePassphrase(formats strfmt.Registry) error {

	if err := validate.Required("passphrase", "body", m.Passphrase); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PassphraseRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PassphraseRequest) UnmarshalBinary(b []byte) error {
	var res PassphraseRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

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

// GenerateAddressesRequest generate addresses request
// swagger:model GenerateAddressesRequest
type GenerateAddressesRequest struct {

	// address n
	// Required: true
	AddressN *int64 `json:"address_n"`

	// confirm address
	ConfirmAddress bool `json:"confirm_address,omitempty"`

	// start index
	StartIndex int64 `json:"start_index,omitempty"`
}

// Validate validates this generate addresses request
func (m *GenerateAddressesRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddressN(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GenerateAddressesRequest) validateAddressN(formats strfmt.Registry) error {

	if err := validate.Required("address_n", "body", m.AddressN); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GenerateAddressesRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GenerateAddressesRequest) UnmarshalBinary(b []byte) error {
	var res GenerateAddressesRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

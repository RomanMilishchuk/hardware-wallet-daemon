// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HTTPSuccessResponse HTTP success response
//
// swagger:model HTTPSuccessResponse
type HTTPSuccessResponse struct {

	// data
	Data []string `json:"data"`
}

// Validate validates this HTTP success response
func (m *HTTPSuccessResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HTTPSuccessResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPSuccessResponse) UnmarshalBinary(b []byte) error {
	var res HTTPSuccessResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

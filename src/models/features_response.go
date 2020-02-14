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

// FeaturesResponse features response
//
// swagger:model FeaturesResponse
type FeaturesResponse struct {

	// data
	Data *FeaturesResponseData `json:"data,omitempty"`
}

// Validate validates this features response
func (m *FeaturesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FeaturesResponse) validateData(formats strfmt.Registry) error {

	if swag.IsZero(m.Data) { // not required
		return nil
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *FeaturesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FeaturesResponse) UnmarshalBinary(b []byte) error {
	var res FeaturesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// FeaturesResponseData features response data
//
// swagger:model FeaturesResponseData
type FeaturesResponseData struct {

	// bootloader hash
	BootloaderHash string `json:"bootloader_hash,omitempty"`

	// device id
	DeviceID string `json:"device_id,omitempty"`

	// firmware features
	// Required: true
	FirmwareFeatures *int64 `json:"firmware_features"`

	// fw major
	// Required: true
	FwMajor *int64 `json:"fw_major"`

	// fw minor
	// Required: true
	FwMinor *int64 `json:"fw_minor"`

	// fw patch
	// Required: true
	FwPatch *int64 `json:"fw_patch"`

	// initialized
	// Required: true
	Initialized *bool `json:"initialized"`

	// label
	Label string `json:"label,omitempty"`

	// major version
	MajorVersion int64 `json:"major_version,omitempty"`

	// minor version
	MinorVersion int64 `json:"minor_version,omitempty"`

	// model
	Model string `json:"model,omitempty"`

	// needs backup
	// Required: true
	NeedsBackup *bool `json:"needs_backup"`

	// passphrase cached
	// Required: true
	PassphraseCached *bool `json:"passphrase_cached"`

	// passphrase protection
	// Required: true
	PassphraseProtection *bool `json:"passphrase_protection"`

	// patch version
	PatchVersion int64 `json:"patch_version,omitempty"`

	// pin cached
	// Required: true
	PinCached *bool `json:"pin_cached"`

	// pin protection
	// Required: true
	PinProtection *bool `json:"pin_protection"`

	// vendor
	// Required: true
	Vendor *string `json:"vendor"`
}

// Validate validates this features response data
func (m *FeaturesResponseData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFirmwareFeatures(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFwMajor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFwMinor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFwPatch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInitialized(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNeedsBackup(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePassphraseCached(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePassphraseProtection(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePinCached(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePinProtection(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVendor(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FeaturesResponseData) validateFirmwareFeatures(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"firmware_features", "body", m.FirmwareFeatures); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateFwMajor(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"fw_major", "body", m.FwMajor); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateFwMinor(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"fw_minor", "body", m.FwMinor); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateFwPatch(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"fw_patch", "body", m.FwPatch); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateInitialized(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"initialized", "body", m.Initialized); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateNeedsBackup(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"needs_backup", "body", m.NeedsBackup); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validatePassphraseCached(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"passphrase_cached", "body", m.PassphraseCached); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validatePassphraseProtection(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"passphrase_protection", "body", m.PassphraseProtection); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validatePinCached(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"pin_cached", "body", m.PinCached); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validatePinProtection(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"pin_protection", "body", m.PinProtection); err != nil {
		return err
	}

	return nil
}

func (m *FeaturesResponseData) validateVendor(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"vendor", "body", m.Vendor); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *FeaturesResponseData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FeaturesResponseData) UnmarshalBinary(b []byte) error {
	var res FeaturesResponseData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

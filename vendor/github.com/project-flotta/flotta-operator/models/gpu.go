// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Gpu gpu
//
// swagger:model gpu
type Gpu struct {

	// Device address (for example "0000:00:02.0")
	Address string `json:"address,omitempty"`

	// ID of the device (for example "3ea0")
	DeviceID string `json:"device_id,omitempty"`

	// Product name of the device (for example "UHD Graphics 620 (Whiskey Lake)")
	Name string `json:"name,omitempty"`

	// The name of the device vendor (for example "Intel Corporation")
	Vendor string `json:"vendor,omitempty"`

	// ID of the vendor (for example "8086")
	VendorID string `json:"vendor_id,omitempty"`
}

// Validate validates this gpu
func (m *Gpu) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this gpu based on context it is used
func (m *Gpu) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Gpu) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Gpu) UnmarshalBinary(b []byte) error {
	var res Gpu
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

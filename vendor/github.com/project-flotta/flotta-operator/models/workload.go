// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Workload workload
//
// swagger:model workload
type Workload struct {

	// Workload Annotations
	Annotations map[string]string `json:"annotations,omitempty"`

	// configmaps
	Configmaps ConfigmapList `json:"configmaps,omitempty"`

	// cron defintion
	Cron string `json:"cron,omitempty"`

	// data
	Data *DataConfiguration `json:"data,omitempty"`

	// image registries
	ImageRegistries *ImageRegistries `json:"imageRegistries,omitempty"`

	// Kind of workload
	Kind string `json:"kind,omitempty"`

	// Workload labels
	Labels map[string]string `json:"labels,omitempty"`

	// Log collection target for this workload
	LogCollection string `json:"log_collection,omitempty"`

	// metrics
	Metrics *Metrics `json:"metrics,omitempty"`

	// Name of the workload
	Name string `json:"name,omitempty"`

	// Namespace of the workload
	Namespace string `json:"namespace,omitempty"`

	// profiles
	Profiles []*WorkloadProfile `json:"profiles"`

	// permission
	Rootless bool `json:"rootless,omitempty"`

	// specification
	Specification string `json:"specification,omitempty"`
}

// Validate validates this workload
func (m *Workload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfigmaps(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImageRegistries(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetrics(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProfiles(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Workload) validateConfigmaps(formats strfmt.Registry) error {
	if swag.IsZero(m.Configmaps) { // not required
		return nil
	}

	if err := m.Configmaps.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("configmaps")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("configmaps")
		}
		return err
	}

	return nil
}

func (m *Workload) validateData(formats strfmt.Registry) error {
	if swag.IsZero(m.Data) { // not required
		return nil
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) validateImageRegistries(formats strfmt.Registry) error {
	if swag.IsZero(m.ImageRegistries) { // not required
		return nil
	}

	if m.ImageRegistries != nil {
		if err := m.ImageRegistries.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("imageRegistries")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("imageRegistries")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) validateMetrics(formats strfmt.Registry) error {
	if swag.IsZero(m.Metrics) { // not required
		return nil
	}

	if m.Metrics != nil {
		if err := m.Metrics.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metrics")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("metrics")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) validateProfiles(formats strfmt.Registry) error {
	if swag.IsZero(m.Profiles) { // not required
		return nil
	}

	for i := 0; i < len(m.Profiles); i++ {
		if swag.IsZero(m.Profiles[i]) { // not required
			continue
		}

		if m.Profiles[i] != nil {
			if err := m.Profiles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("profiles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("profiles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this workload based on the context it is used
func (m *Workload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateConfigmaps(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateImageRegistries(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMetrics(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProfiles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Workload) contextValidateConfigmaps(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Configmaps.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("configmaps")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("configmaps")
		}
		return err
	}

	return nil
}

func (m *Workload) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if m.Data != nil {
		if err := m.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) contextValidateImageRegistries(ctx context.Context, formats strfmt.Registry) error {

	if m.ImageRegistries != nil {
		if err := m.ImageRegistries.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("imageRegistries")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("imageRegistries")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) contextValidateMetrics(ctx context.Context, formats strfmt.Registry) error {

	if m.Metrics != nil {
		if err := m.Metrics.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metrics")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("metrics")
			}
			return err
		}
	}

	return nil
}

func (m *Workload) contextValidateProfiles(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Profiles); i++ {

		if m.Profiles[i] != nil {
			if err := m.Profiles[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("profiles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("profiles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Workload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Workload) UnmarshalBinary(b []byte) error {
	var res Workload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

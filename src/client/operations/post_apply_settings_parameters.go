// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/SkycoinProject/hardware-wallet-daemon/src/models"
)

// NewPostApplySettingsParams creates a new PostApplySettingsParams object
// with the default values initialized.
func NewPostApplySettingsParams() *PostApplySettingsParams {
	var ()
	return &PostApplySettingsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostApplySettingsParamsWithTimeout creates a new PostApplySettingsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostApplySettingsParamsWithTimeout(timeout time.Duration) *PostApplySettingsParams {
	var ()
	return &PostApplySettingsParams{

		timeout: timeout,
	}
}

// NewPostApplySettingsParamsWithContext creates a new PostApplySettingsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostApplySettingsParamsWithContext(ctx context.Context) *PostApplySettingsParams {
	var ()
	return &PostApplySettingsParams{

		Context: ctx,
	}
}

// NewPostApplySettingsParamsWithHTTPClient creates a new PostApplySettingsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostApplySettingsParamsWithHTTPClient(client *http.Client) *PostApplySettingsParams {
	var ()
	return &PostApplySettingsParams{
		HTTPClient: client,
	}
}

/*PostApplySettingsParams contains all the parameters to send to the API endpoint
for the post apply settings operation typically these are written to a http.Request
*/
type PostApplySettingsParams struct {

	/*ApplySettingsRequest
	  ApplySettingsRequest is request data for /api/v1/apply_settings

	*/
	ApplySettingsRequest *models.ApplySettingsRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post apply settings params
func (o *PostApplySettingsParams) WithTimeout(timeout time.Duration) *PostApplySettingsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post apply settings params
func (o *PostApplySettingsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post apply settings params
func (o *PostApplySettingsParams) WithContext(ctx context.Context) *PostApplySettingsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post apply settings params
func (o *PostApplySettingsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post apply settings params
func (o *PostApplySettingsParams) WithHTTPClient(client *http.Client) *PostApplySettingsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post apply settings params
func (o *PostApplySettingsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApplySettingsRequest adds the applySettingsRequest to the post apply settings params
func (o *PostApplySettingsParams) WithApplySettingsRequest(applySettingsRequest *models.ApplySettingsRequest) *PostApplySettingsParams {
	o.SetApplySettingsRequest(applySettingsRequest)
	return o
}

// SetApplySettingsRequest adds the applySettingsRequest to the post apply settings params
func (o *PostApplySettingsParams) SetApplySettingsRequest(applySettingsRequest *models.ApplySettingsRequest) {
	o.ApplySettingsRequest = applySettingsRequest
}

// WriteToRequest writes these params to a swagger request
func (o *PostApplySettingsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ApplySettingsRequest != nil {
		if err := r.SetBodyParam(o.ApplySettingsRequest); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

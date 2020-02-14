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
)

// NewPutFirmwareUpdateParams creates a new PutFirmwareUpdateParams object
// with the default values initialized.
func NewPutFirmwareUpdateParams() *PutFirmwareUpdateParams {

	return &PutFirmwareUpdateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPutFirmwareUpdateParamsWithTimeout creates a new PutFirmwareUpdateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPutFirmwareUpdateParamsWithTimeout(timeout time.Duration) *PutFirmwareUpdateParams {

	return &PutFirmwareUpdateParams{

		timeout: timeout,
	}
}

// NewPutFirmwareUpdateParamsWithContext creates a new PutFirmwareUpdateParams object
// with the default values initialized, and the ability to set a context for a request
func NewPutFirmwareUpdateParamsWithContext(ctx context.Context) *PutFirmwareUpdateParams {

	return &PutFirmwareUpdateParams{

		Context: ctx,
	}
}

// NewPutFirmwareUpdateParamsWithHTTPClient creates a new PutFirmwareUpdateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPutFirmwareUpdateParamsWithHTTPClient(client *http.Client) *PutFirmwareUpdateParams {

	return &PutFirmwareUpdateParams{
		HTTPClient: client,
	}
}

/*PutFirmwareUpdateParams contains all the parameters to send to the API endpoint
for the put firmware update operation typically these are written to a http.Request
*/
type PutFirmwareUpdateParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the put firmware update params
func (o *PutFirmwareUpdateParams) WithTimeout(timeout time.Duration) *PutFirmwareUpdateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put firmware update params
func (o *PutFirmwareUpdateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put firmware update params
func (o *PutFirmwareUpdateParams) WithContext(ctx context.Context) *PutFirmwareUpdateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put firmware update params
func (o *PutFirmwareUpdateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put firmware update params
func (o *PutFirmwareUpdateParams) WithHTTPClient(client *http.Client) *PutFirmwareUpdateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put firmware update params
func (o *PutFirmwareUpdateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *PutFirmwareUpdateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

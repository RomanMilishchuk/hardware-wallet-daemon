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

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetCsrfParams creates a new GetCsrfParams object
// with the default values initialized.
func NewGetCsrfParams() *GetCsrfParams {

	return &GetCsrfParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCsrfParamsWithTimeout creates a new GetCsrfParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCsrfParamsWithTimeout(timeout time.Duration) *GetCsrfParams {

	return &GetCsrfParams{

		timeout: timeout,
	}
}

// NewGetCsrfParamsWithContext creates a new GetCsrfParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCsrfParamsWithContext(ctx context.Context) *GetCsrfParams {

	return &GetCsrfParams{

		Context: ctx,
	}
}

// NewGetCsrfParamsWithHTTPClient creates a new GetCsrfParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCsrfParamsWithHTTPClient(client *http.Client) *GetCsrfParams {

	return &GetCsrfParams{
		HTTPClient: client,
	}
}

/*GetCsrfParams contains all the parameters to send to the API endpoint
for the get csrf operation typically these are written to a http.Request
*/
type GetCsrfParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get csrf params
func (o *GetCsrfParams) WithTimeout(timeout time.Duration) *GetCsrfParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get csrf params
func (o *GetCsrfParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get csrf params
func (o *GetCsrfParams) WithContext(ctx context.Context) *GetCsrfParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get csrf params
func (o *GetCsrfParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get csrf params
func (o *GetCsrfParams) WithHTTPClient(client *http.Client) *GetCsrfParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get csrf params
func (o *GetCsrfParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetCsrfParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
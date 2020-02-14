// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/SkycoinProject/hardware-wallet-daemon/src/models"
)

// GetCsrfReader is a Reader for the GetCsrf structure.
type GetCsrfReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCsrfReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetCsrfOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetCsrfDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCsrfOK creates a GetCsrfOK with default headers values
func NewGetCsrfOK() *GetCsrfOK {
	return &GetCsrfOK{}
}

/*GetCsrfOK handles this case with default header values.

successful operation
*/
type GetCsrfOK struct {
	Payload *models.CSRFResponse
}

func (o *GetCsrfOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CSRFResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCsrfDefault creates a GetCsrfDefault with default headers values
func NewGetCsrfDefault(code int) *GetCsrfDefault {
	return &GetCsrfDefault{
		_statusCode: code,
	}
}

/*GetCsrfDefault handles this case with default header values.

error
*/
type GetCsrfDefault struct {
	_statusCode int

	Payload *models.HTTPErrorResponse
}

// Code gets the status code for the get csrf default response
func (o *GetCsrfDefault) Code() int {
	return o._statusCode
}

func (o *GetCsrfDefault) Error() string {
	return o.Payload.Error.Message
}

func (o *GetCsrfDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

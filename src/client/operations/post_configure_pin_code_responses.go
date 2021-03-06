// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/skycoin/hardware-wallet-daemon/src/models"
)

// PostConfigurePinCodeReader is a Reader for the PostConfigurePinCode structure.
type PostConfigurePinCodeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostConfigurePinCodeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPostConfigurePinCodeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewPostConfigurePinCodeDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostConfigurePinCodeOK creates a PostConfigurePinCodeOK with default headers values
func NewPostConfigurePinCodeOK() *PostConfigurePinCodeOK {
	return &PostConfigurePinCodeOK{}
}

/*PostConfigurePinCodeOK handles this case with default header values.

success
*/
type PostConfigurePinCodeOK struct {
	Payload *models.HttpsuccessResponse
}

func (o *PostConfigurePinCodeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttpsuccessResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostConfigurePinCodeDefault creates a PostConfigurePinCodeDefault with default headers values
func NewPostConfigurePinCodeDefault(code int) *PostConfigurePinCodeDefault {
	return &PostConfigurePinCodeDefault{
		_statusCode: code,
	}
}

/*PostConfigurePinCodeDefault handles this case with default header values.

error
*/
type PostConfigurePinCodeDefault struct {
	_statusCode int

	Payload *models.HTTPErrorResponse
}

// Code gets the status code for the post configure pin code default response
func (o *PostConfigurePinCodeDefault) Code() int {
	return o._statusCode
}

func (o *PostConfigurePinCodeDefault) Error() string {
	return o.Payload.Error.Message
}

func (o *PostConfigurePinCodeDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

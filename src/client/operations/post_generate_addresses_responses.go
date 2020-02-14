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

// PostGenerateAddressesReader is a Reader for the PostGenerateAddresses structure.
type PostGenerateAddressesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostGenerateAddressesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPostGenerateAddressesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewPostGenerateAddressesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostGenerateAddressesOK creates a PostGenerateAddressesOK with default headers values
func NewPostGenerateAddressesOK() *PostGenerateAddressesOK {
	return &PostGenerateAddressesOK{}
}

/*PostGenerateAddressesOK handles this case with default header values.

success
*/
type PostGenerateAddressesOK struct {
	Payload *models.GenerateAddressesResponse
}

func (o *PostGenerateAddressesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenerateAddressesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostGenerateAddressesDefault creates a PostGenerateAddressesDefault with default headers values
func NewPostGenerateAddressesDefault(code int) *PostGenerateAddressesDefault {
	return &PostGenerateAddressesDefault{
		_statusCode: code,
	}
}

/*PostGenerateAddressesDefault handles this case with default header values.

error
*/
type PostGenerateAddressesDefault struct {
	_statusCode int

	Payload *models.HTTPErrorResponse
}

// Code gets the status code for the post generate addresses default response
func (o *PostGenerateAddressesDefault) Code() int {
	return o._statusCode
}

func (o *PostGenerateAddressesDefault) Error() string {
	return o.Payload.Error.Message
}

func (o *PostGenerateAddressesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

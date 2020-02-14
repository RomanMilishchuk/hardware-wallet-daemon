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

// PostBackupReader is a Reader for the PostBackup structure.
type PostBackupReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostBackupReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPostBackupOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewPostBackupDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostBackupOK creates a PostBackupOK with default headers values
func NewPostBackupOK() *PostBackupOK {
	return &PostBackupOK{}
}

/*PostBackupOK handles this case with default header values.

success
*/
type PostBackupOK struct {
	Payload *models.HTTPSuccessResponse
}

func (o *PostBackupOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HTTPSuccessResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostBackupDefault creates a PostBackupDefault with default headers values
func NewPostBackupDefault(code int) *PostBackupDefault {
	return &PostBackupDefault{
		_statusCode: code,
	}
}

/*PostBackupDefault handles this case with default header values.

error
*/
type PostBackupDefault struct {
	_statusCode int

	Payload *models.HTTPErrorResponse
}

// Code gets the status code for the post backup default response
func (o *PostBackupDefault) Code() int {
	return o._statusCode
}

func (o *PostBackupDefault) Error() string {
	return o.Payload.Error.Message
}

func (o *PostBackupDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

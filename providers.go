package requist

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io"
	"strings"
)

//=== Request Body manipulators

// BodyProvider provides Body content for http.Request attachment.
type BodyProvider interface {
	// ContentType returns the Content-Type of the body.
	ContentType() string
	// Body returns the io.Reader body.
	Body() (io.Reader, error)
}

//=== FormProvider implementation of BodyProvider interface

// formProvider implementation of BodyProvider interface
type formProvider struct {
	payload interface{}
}

// formProvider ContentType just returns FormContentType for validations
func (p formProvider) ContentType() string {

	return FormContentType
}

// formProvider Body prepare our request body in Forms format
func (p formProvider) Body() (io.Reader, error) {

	values, err := query.Values(p.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}

//=== JsonProvider implementation of BodyProvider interface

// jsonProvider implementation of BodyProvider interface
type jsonProvider struct {
	payload interface{}
}

// jsonProvider ContentType just returns JsonContentType for validations
func (p jsonProvider) ContentType() string {

	return JsonContentType
}

// jsonProvider Body prepare our request body in JSON format
func (p jsonProvider) Body() (io.Reader, error) {

	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(p.payload)

	if err != nil {
		return nil, err
	}
	return buffer, nil
}

//=== Plain Text Provider implementation of BodyProvider interface

// textProvider implementation of BodyProvider interface
type textProvider struct {
	payload interface{}
}

// formProvider ContentType just returns FormContentType for validations
func (p textProvider) ContentType() string {

	return TextContentType
}

// textProvider Body prepare our request body in Forms format
func (p textProvider) Body() (io.Reader, error) {

	return nil, nil
}

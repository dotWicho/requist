package requist

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//=== Response Body manipulators

// BodyResponse decodes http responses into struct values.
type BodyResponse interface {
	// Decode decodes the response into the value pointed to by v.
	Accept() string
	Decode(resp io.Reader, v interface{}) (err error)
}

// formResponse decodes http response FORM into a map[string]string.
type formResponse struct{}

// Accept just return the Accept Type (application/x-www-form-urlencoded)
func (r formResponse) Accept() string {
	return FormContentType
}

// Decode decodes the Response Body into the value pointed to by v
// 	Must be provided a non-nil forms (struct) reference
func (r formResponse) Decode(resp io.Reader, v interface{}) (err error) {

	if v, err = ioutil.ReadAll(resp); err != nil {
		return err
	}
	return nil
}

// jsonResponse decodes http response JSON into a JSON-tagged struct value.
type jsonResponse struct{}

// Accept just return the Accept Type (application/json)
func (r jsonResponse) Accept() string {
	return JsonContentType
}

// Decode decodes the Response Body into the value pointed to by v
// 	Must be provided a non-nil json (struct) reference
func (r jsonResponse) Decode(resp io.Reader, v interface{}) (err error) {

	if err = json.NewDecoder(resp).Decode(v); err != nil {
		return err
	}
	return nil
}

// textResponse decodes http response into a simple plain text.
type textResponse struct{}

// Accept just return the Accept Type (text/plain)
func (r textResponse) Accept() string {
	return TextContentType
}

// Decode decodes the Response Body into the value pointed to by v.
// 	Must be provided a non-nil plain text (string) reference
func (r textResponse) Decode(resp io.Reader, v interface{}) (err error) {

	if v, err = ioutil.ReadAll(resp); err != nil {
		return err
	}
	return nil
}

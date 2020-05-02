package requist

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	// "fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	// "reflect"
	"strings"
	"time"

	// Easy handle for QueryStrings
	"github.com/google/go-querystring/query"
)

//#$$=== Useful constants

const (
	acceptHeader    string = "Accept"
	contentType     string = "Content-Type"
	textContentType string = "text/plain"
	jsonContentType string = "application/json"
	formContentType string = "application/x-www-form-urlencoded"
)

//#$$=== Request Body manipulators

// BodyProvider provides Body content for http.Request attachment.
type BodyProvider interface {
	// ContentType returns the Content-Type of the body.
	ContentType() string
	// Body returns the io.Reader body.
	Body() (io.Reader, error)
}

// bodyProvider provides the wrapped body value as a Body for reqests.
type bodyProvider struct {
	body io.Reader
}

// bodyProvider ContentType() just return a empty Content-Type
func (p bodyProvider) ContentType() string {

	return ""
}

// bodyProvider Body() returns the body request
func (p bodyProvider) Body() (io.Reader, error) {

	return p.body, nil
}

//#$$=== FormProvider implementation of BodyProvider interface

type formProvider struct {
	payload interface{}
}

// formProvider ContentType just returns formContentType for validations
func (p formProvider) ContentType() string {

	return formContentType
}

// formProvider Body prepare our request body in Forms format
func (p formProvider) Body() (io.Reader, error) {

	values, err := query.Values(p.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}

//#$$=== JsonProvider implementation of BodyProvider interface

type jsonProvider struct {
	payload interface{}
}

// jsonProvider ContentType just returns jsonContentType for validations
func (p jsonProvider) ContentType() string {

	return jsonContentType
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

//#$$=== FormProvider implementation of BodyProvider interface

type textProvider struct {
	payload interface{}
}

// formProvider ContentType just returns formContentType for validations
func (p textProvider) ContentType() string {

	return textContentType
}

// textProvider Body prepare our request body in Forms format
func (p textProvider) Body() (io.Reader, error) {

	return nil, nil
}

//#$$=== Response Body manipulators

// BodyResponse decodes http responses into struct values.
type BodyResponse interface {
	// Decode decodes the response into the value pointed to by v.
	Accept() string
	Decode(resp io.Reader, v interface{}) error
}

// bodyResponse provides the wrapped to handle response body from requests.
type bodyResponse struct{}

// formResponse decodes http response FORM into a map[string]string.
type formResponse struct{}

// Accept just return the Accept Type (application/x-www-form-urlencoded)
func (r formResponse) Accept() string {
	return formContentType
}

func (r formResponse) Decode(resp io.Reader, v interface{}) error {

	contents, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}
	v = string(contents)
	return nil
}

// jsonResponse decodes http response JSON into a JSON-tagged struct value.
type jsonResponse struct{}

// Accept just return the Accept Type (application/json)
func (r jsonResponse) Accept() string {
	return jsonContentType
}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (r jsonResponse) Decode(resp io.Reader, v interface{}) error {

	if err := json.NewDecoder(resp).Decode(v); err != nil {
		return err
	}

	return nil
}

// textResponse decodes http response into a simple plain text.
type textResponse struct{}

// Accept just return the Accept Type (text/plain)
func (r textResponse) Accept() string {
	return textContentType
}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (r textResponse) Decode(resp io.Reader, v interface{}) error {

	contents, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}
	v = string(contents)
	return nil
}

//#$$=== Definitions: struct for Requests manipulations

// Requist struct Encapsulate an HTTP(S) requests builder and sender
type Requist struct {
	auth   string
	method string
	url    string
	path   string

	statuscode int

	client   *http.Client
	header   *http.Header
	querys   *url.Values
	provider BodyProvider
	response BodyResponse
}

//#$$=== Supplemental functions to manipulate path

// ParseBaseURL check if is valid the base string pased
func parseBaseURL(base string) string {
	url, err := url.Parse(base)
	if err != nil {
		log.Fatalln()
	}
	url.RawQuery = ""
	url.Fragment = ""

	return url.String()
}

// ParsePathURL check relative path
func parsePathURL(base string, path string) string {
	url, err := url.Parse(base)
	if err != nil {
		log.Fatalln()
	}
	url.Path = path
	return url.String()
}

//#$$=== Functions to create a Requist instance

// New function
//  @param baseURL
//  @return Requist class pointer
//
func New(baseURL string) *Requist {
	requist := new(Requist)

	if baseURL != "" {
		requist.url = parseBaseURL(baseURL)
	}
	requist.header = &http.Header{}
	requist.querys = &url.Values{}
	requist.client = &http.Client{}
	requist.client.Timeout = 60 * time.Second
	requist.provider = nil
	requist.response = nil

	return requist.Base(requist.url)
}

// New class function
//  @param Requist class pointer, previous existing instance
//  @return Requist class pointer who clone some data from passed class
//
func (r *Requist) New(baseURL string) *Requist {
	requist := new(Requist)

	if baseURL != "" {
		requist.url = parseBaseURL(baseURL)
	}
	requist.header = r.header
	requist.querys = r.querys
	requist.provider = r.provider
	requist.response = r.response

	return requist.Base(requist.url)
}

//#$$=== Core function of Requist class

// Request ... Here it's where the magic show up
func (r *Requist) Request(successV, failureV interface{}) (*Requist, error) {

	requestPath, err := r.addQueryParams(r.url)
	if err != nil {
		return r, err
	}

	var body io.Reader
	if r.provider != nil {

		body, err = r.provider.Body()
		if err != nil {
			return r, err
		}
	}

	// Prepares request struct with all fields needed
	request, err := http.NewRequest(r.method, requestPath, body)
	if err != nil {
		return r, err
	}
	// Proceed to clone headers pre populated to the request class
	request.Header = r.header.Clone()

	// Fire up the request agains the server
	response, err := r.client.Do(request)
	if err != nil {
		return r, err
	}
	// Defer close response body
	defer response.Body.Close()
	defer r.CleanQueryParams()

	r.statuscode = response.StatusCode

	// Decode from r.response Accept() type
	if (successV != nil || failureV != nil) && r.statuscode != 204 {
		if 200 <= r.statuscode && r.statuscode <= 299 {
			if successV != nil {

				if r.response != nil {
					if err := r.response.Decode(response.Body, successV); err != nil {
						return r, err
					}
				}
			}
		} else {
			if failureV != nil {

				if r.response != nil {
					if err := r.response.Decode(response.Body, failureV); err != nil {
						return r, err
					}
				}
			}
		}
	}
	return r, err
}

//#$$=== Provider Body functions, used to set de content of payload to send on request

// BodyProvider sets the Requests's body provider from original BodyProvider interface{}.
func (r *Requist) BodyProvider(body BodyProvider) *Requist {

	if body == nil {
		return r
	}

	ct := body.ContentType()
	if ct != "" {
		r.provider = body
		r.SetHeader(contentType, ct)
		r.Accept(ct)
	}

	return r
}

// BodyAsForm sets the Requests's body from a formProvider
func (r *Requist) BodyAsForm(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(formProvider{payload: body})
}

// BodyAsJSON sets the Requests's body from a jsonProvider
func (r *Requist) BodyAsJSON(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(jsonProvider{payload: body})
}

// BodyAsForm sets the Requests's body from a formProvider
func (r *Requist) BodyAsText(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(textProvider{payload: body})
}

//#$$===

func (r *Requist) BodyResponse(body BodyResponse) *Requist {

	if body == nil {
		return r
	}

	ct := body.Accept()
	if ct != "" {
		r.response = body
		r.SetHeader(acceptHeader, ct)
	}

	return r
}

func (r *Requist) Accept(accept string) {
	switch accept {
	case formContentType:
		r.BodyResponse(formResponse{})
	case jsonContentType:
		r.BodyResponse(jsonResponse{})
	case textContentType:
		r.BodyResponse(textResponse{})
	default:
		r.response = nil
	}
}

//#$$=== QueryParams manipulation functions

// addQueryParams ...
func (r *Requist) addQueryParams(basePath string) (string, error) {
	reqURL, err := url.Parse(parseBaseURL(r.url))

	reqURL.RawQuery = r.querys.Encode()

	return reqURL.String(), err
}

//#$$=== Header manipulation functions

// Add adds the key, value pair in Headers, appending values for existing keys
// to the key's values. Header keys are canonicalized.
func (r *Requist) AddHeader(key, value string) *Requist {

	r.header.Add(key, value)

	return r
}

// Set sets the key, value pair in Headers, replacing existing values
// associated with key. Header keys are canonicalized.
func (r *Requist) SetHeader(key, value string) *Requist {

	r.header.Set(key, value)

	return r
}

// Remove the key, value pair in Headers
func (r *Requist) DelHeader(key, value string) *Requist {

	r.header.Del(key)

	return r
}

// AddQueryParam adds the key, value tuples in QueryParams, appending values for existing keys
func (r *Requist) AddQueryParam(key, value string) *Requist {

	if r.querys != nil {
		r.querys.Add(key, value)
	}
	return r
}

// SetQueryParam set the key, value tuples in params to
func (r *Requist) SetQueryParam(key, value string) *Requist {

	r.querys.Set(key, value)

	return r
}

// Remove the key from QueryParams
func (r *Requist) DelQueryParam(key, value string) *Requist {

	r.querys.Del(key)

	return r
}

// Remove the key from QueryParams
func (r *Requist) CleanQueryParams() *Requist {

	r.querys = &url.Values{}
	return r
}

// SetBasicAuth sets the Authorization header to use HTTP Basic Authentication
func (r *Requist) SetBasicAuth(username, password string) *Requist {

	if username != "" && password != "" {
		r.auth = username + ":" + password
		r.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(r.auth)))
	}

	return r
}

//=== Utilities functions, used to return some values from Requist class

// StatusCode return the HTTP StatusCode from last request
func (r *Requist) StatusCode() int {

	return r.statuscode
}

// GetBasicAuth return the auth stored at the Requist class
func (r *Requist) GetBasicAuth() string {

	return r.auth
}

//=== Utilities functions, to set up URL base, URL path, HTTP method to use...

// Base sets base url to use for a client
func (r *Requist) Base(base string) *Requist {

	r.url = parseBaseURL(base)

	return r
}

// Path sets request path to use in next request
func (r *Requist) Path(path string) *Requist {

	r.url = parsePathURL(r.url, path)

	return r
}

// Method set HTTP Method to execute
func (r *Requist) Method(method string) *Requist {

	r.method = method
	return r
}

//#$$=== Requist functions executers, Correspond to HTTP Methods

// This implement HEAD HTTP Method
func (r *Requist) Head(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodHead).Path(path).Request(successV, failureV)
}

// This implement GET HTTP Method
func (r *Requist) Get(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodGet).Path(path).Request(successV, failureV)
}

// This implement PUT HTTP Method
func (r *Requist) Put(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodPut).Path(path).Request(successV, failureV)
}

// This implement POST HTTP Method
func (r *Requist) Post(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodPost).Path(path).Request(successV, failureV)
}

// This implement PATCH HTTP Method
func (r *Requist) Patch(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodPatch).Path(path).Request(successV, failureV)
}

// This implement DELETE HTTP Method
func (r *Requist) Delete(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodDelete).Path(path).Request(successV, failureV)
}

// This implement OPTIONS HTTP Method
func (r *Requist) Options(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodOptions).Path(path).Request(successV, failureV)
}

// This implement TRACE HTTP Method
func (r *Requist) Trace(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodTrace).Path(path).Request(successV, failureV)
}

// This implement CONNECT HTTP Method
func (r *Requist) Connect(path string, successV, failureV interface{}) (*Requist, error) {

	return r.Method(http.MethodConnect).Path(path).Request(successV, failureV)
}

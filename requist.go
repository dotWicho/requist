package requist

import (
	"encoding/base64"
	// We use go-cleanhttp because it contains a better implementation of http.Transport
	// and allows us to abstract from these changes
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"net/http"
	"net/url"
	"time"
)

//=== Requests manipulations interface

// Requist interface Define all Methods
type requist interface {
	New(baseURL string) *Requist

	SetClientTransport(transport *http.Transport)
	SetClientTimeout(timeout time.Duration)
	BodyProvider(body BodyProvider) *Requist
	BodyAsForm(body interface{}) *Requist
	BodyAsJSON(body interface{}) *Requist
	BodyAsText(body interface{}) *Requist
	BodyResponse(body BodyResponse) *Requist
	Accept(accept string)

	addQueryParams() (string, error)

	AddHeader(key, value string)
	SetHeader(key, value string)
	DelHeader(key string)
	AddQueryParam(key, value string)
	SetQueryParam(key, value string)
	DelQueryParam(key string)
	CleanQueryParams()
	SetBasicAuth(username, password string) *Requist
	StatusCode() int
	GetBasicAuth() string
	Base(base string) *Requist
	Path(path string) *Requist
	Method(method string) *Requist
	Head(path string, success, failure interface{}) (*Requist, error)
	Get(path string, success, failure interface{}) (*Requist, error)
	Put(path string, success, failure interface{}) (*Requist, error)
	Post(path string, success, failure interface{}) (*Requist, error)
	Patch(path string, success, failure interface{}) (*Requist, error)
	Delete(path string, success, failure interface{}) (*Requist, error)
	Options(path string, success, failure interface{}) (*Requist, error)
	Trace(path string, success, failure interface{}) (*Requist, error)
	Connect(path string, success, failure interface{}) (*Requist, error)
}

// Requist struct Encapsulate an HTTP(S) requests builder and sender
type Requist struct {

	// Basics
	auth   string
	method string
	url    string
	path   string

	// Holds last HTTP Response Code
	statuscode int

	// Handle HTTP(S) primitives
	client  *http.Client
	header  *http.Header
	queries *url.Values

	// Bodies, Request and Response
	provider BodyProvider
	response BodyResponse
}

//=== Functions to create a Requist instance

// New function
//  @param baseURL
//  @return Requist class pointer
//
func New(baseURL string) *Requist {

	r := new(Requist)

	if baseURL != "" {
		r.url = parseBaseURL(baseURL)
	}
	r.header = &http.Header{}
	r.queries = &url.Values{}
	r.client = &http.Client{}
	r.SetClientTransport(cleanhttp.DefaultTransport())
	r.SetClientTimeout(defaultTimeout)
	r.provider = nil
	r.response = nil

	return r.Base(r.url)
}

// New class function
//  @param Requist class pointer, previous existing instance
//  @return Requist class pointer who clone some data from passed class
//
func (r *Requist) New(baseURL string) *Requist {
	req := new(Requist)

	if baseURL != "" {
		req.url = parseBaseURL(baseURL)
	}
	req.header = r.header
	req.queries = r.queries
	req.provider = r.provider
	req.response = r.response

	return req.Base(req.url)
}

// SetClientTransport take transport param and set client HTTP Transport
func (r *Requist) SetClientTransport(transport *http.Transport) {

	r.client.Transport = transport
}

// SetClientTimeout take timeout param and set client Timeout seconds based
func (r *Requist) SetClientTimeout(timeout time.Duration) {

	r.client.Timeout = timeout * time.Second
}

//#$$=== Core function of Requist class

// Request ... Here it's where the magic show up
func (r *Requist) Request(success, failure interface{}) (*Requist, error) {

	requestPath, err := r.addQueryParams()
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
	if (success != nil || failure != nil) && r.statuscode != 204 {
		if 200 <= r.statuscode && r.statuscode <= 299 {
			if success != nil {

				if r.response != nil {
					if err := r.response.Decode(response.Body, success); err != nil {
						return r, err
					}
				}
			}
		} else {
			if failure != nil {

				if r.response != nil {
					if err := r.response.Decode(response.Body, failure); err != nil {
						return r, err
					}
				}
			}
		}
	}
	return r, err
}

//#$$=== Provider Body functions, used to set type of payload send on request

// BodyProvider sets the Requests's body provider from original BodyProvider interface{}
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

// BodyAsText sets the Requests's body from a textProvider
func (r *Requist) BodyAsText(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(textProvider{payload: body})
}

//#$$=== Response Body functions, used to set type of response

// BodyResponse sets the response's body
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

// Accept sets the response's body mimeType
func (r *Requist) Accept(accept string) {
	switch accept {
	case FormContentType:
		r.BodyResponse(formResponse{})
	case JsonContentType:
		r.BodyResponse(jsonResponse{})
	case TextContentType:
		r.BodyResponse(textResponse{})
	default:
		r.response = nil
	}
}

//#$$=== QueryParams manipulation functions

// addQueryParams ...
func (r *Requist) addQueryParams() (string, error) {

	reqURL, err := url.Parse(parseBaseURL(r.url))
	if err != nil {
		return "", err
	}
	reqURL.RawQuery = r.queries.Encode()

	return reqURL.String(), err
}

//#$$=== Header manipulation functions

// AddHeader adds the key, value pair in Headers, appending values for existing keys
// to the key's values. Header keys are canonicalized.
func (r *Requist) AddHeader(key, value string) {

	r.header.Add(key, value)
}

// SetHeader sets the key, value pair in Headers, replacing existing values
// associated with key. Header keys are canonicalized.
func (r *Requist) SetHeader(key, value string) {

	r.header.Set(key, value)
}

// DelHeader remove the key, value pair in Headers
func (r *Requist) DelHeader(key string) {

	r.header.Del(key)
}

// AddQueryParam adds the key, value tuples in QueryParams, appending values for existing keys
func (r *Requist) AddQueryParam(key, value string) {

	if r.queries != nil {
		r.queries.Add(key, value)
	}
}

// SetQueryParam set the key, value tuples in params to
func (r *Requist) SetQueryParam(key, value string) {

	if r.queries != nil {
		r.queries.Set(key, value)
	}
}

// DelQueryParam remove the key from QueryParams
func (r *Requist) DelQueryParam(key string) {

	if r.queries != nil {
		r.queries.Del(key)
	}
}

// CleanQueryParams remove all keys from QueryParams
func (r *Requist) CleanQueryParams() {

	r.queries = &url.Values{}
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

// Head implement HEAD HTTP Method
func (r *Requist) Head(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodHead).Path(path).Request(success, failure)
}

// Get implement GET HTTP Method
func (r *Requist) Get(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodGet).Path(path).Request(success, failure)
}

// Put implement PUT HTTP Method
func (r *Requist) Put(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodPut).Path(path).Request(success, failure)
}

// Post implement POST HTTP Method
func (r *Requist) Post(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodPost).Path(path).Request(success, failure)
}

// Patch implement PATCH HTTP Method
func (r *Requist) Patch(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodPatch).Path(path).Request(success, failure)
}

// Delete implement DELETE HTTP Method
func (r *Requist) Delete(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodDelete).Path(path).Request(success, failure)
}

// Options implement OPTIONS HTTP Method
func (r *Requist) Options(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodOptions).Path(path).Request(success, failure)
}

// Trace implement TRACE HTTP Method
func (r *Requist) Trace(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodTrace).Path(path).Request(success, failure)
}

// Connect implement CONNECT HTTP Method
func (r *Requist) Connect(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodConnect).Path(path).Request(success, failure)
}

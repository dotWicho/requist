package requist

import (
	"context"
	"encoding/base64"
	"github.com/dotWicho/logger"
	"strings"

	// We use go-cleanhttp because it contains a better implementation of http.Transport
	// and allows us to abstract from these changes
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Logger is the default logs handler
var Logger = logger.NewLogger(false)

//=== Requests manipulations interface

// Operations interface Define all Methods
type Operations interface {
	SetClientTransport(transport *http.Transport)
	SetClientTimeout(timeout time.Duration)
	SetClientContext(context context.Context)

	BodyProvider(body BodyProvider) *Requist
	BodyAsForm(body interface{}) *Requist
	BodyAsJSON(body interface{}) *Requist
	BodyAsText(body interface{}) *Requist
	BodyResponse(body BodyResponse) *Requist
	Accept(accept string)

	PrepareRequestURI() (string, error)

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
	URI(uri string) *Requist
	Method(method string) *Requist

	Get(path string, success, failure interface{}) (*Requist, error)
	Put(path string, success, failure interface{}) (*Requist, error)
	Post(path string, success, failure interface{}) (*Requist, error)
	Patch(path string, success, failure interface{}) (*Requist, error)
	Delete(path string, success, failure interface{}) (*Requist, error)
	Options(path string, success, failure interface{}) (*Requist, error)
	Connect(path string, success, failure interface{}) (*Requist, error)
}

// Requist struct Encapsulate an HTTP(S) requests builder and sender
type Requist struct {

	// Basics
	auth   string
	method string
	url    string
	uri    string
	path   string

	// Holds last HTTP Response Code
	statuscode int

	// Handle HTTP(S) primitives
	client  *http.Client
	header  *http.Header
	queries *url.Values
	ctx     context.Context

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

	Logger.Debug("Creating a new Client")
	r := &Requist{}

	r.url = ParseBaseURL(baseURL)
	if r.url == "" {
		Logger.Error("Invalid baseURL = %s", baseURL)
		return nil
	}
	r.header = &http.Header{}
	r.queries = &url.Values{}
	r.client = &http.Client{}
	r.ctx = context.Background()
	r.SetClientTransport(cleanhttp.DefaultTransport())
	r.SetClientTimeout(defaultTimeout)
	r.provider = nil
	r.response = nil

	return r.Base(r.url)
}

// SetClientTransport take transport param and set client HTTP Transport
func (r *Requist) SetClientTransport(transport *http.Transport) {

	Logger.Debug("Setting Client Transport %+v", transport)

	r.client.Transport = transport
}

// SetClientTimeout take timeout param and set client Timeout seconds based
func (r *Requist) SetClientTimeout(timeout time.Duration) {

	Logger.Debug("Setting Client Timeout %+v", timeout)

	r.client.Timeout = timeout
}

// SetClientContext take transport param and set client HTTP Transport
func (r *Requist) SetClientContext(context context.Context) {

	Logger.Debug("Setting Client Context %+v", context)
	r.ctx = context
}

//#$$=== Core function of Requist class

// Request ... Here it's where the magic show up
func (r *Requist) Request(success, failure interface{}) (*Requist, error) {

	Logger.Debug("Firing a Request %T %T", success, failure)

	var requestPath string
	var err error

	if requestPath, err = r.PrepareRequestURI(); err != nil {
		return r, err
	}
	Logger.Debug("Request URI to %s", requestPath)

	var body io.Reader
	if r.provider != nil {

		body, err = r.provider.Body()
		if err != nil {
			return r, err
		}
	}

	// Prepares request struct with all fields needed
	var request *http.Request

	if request, err = http.NewRequestWithContext(r.ctx, r.method, requestPath, body); err != nil {
		return r, err
	}

	// Proceed to clone headers pre populated to the request class
	request.Header = r.header.Clone()

	// Fire up the request against the server
	var response *http.Response
	if response, err = r.client.Do(request); err != nil {
		return r, err
	}

	// Defer close response body
	defer response.Body.Close()
	defer r.CleanQueryParams()

	// backup response StatusCode into Requist.statuscode
	r.statuscode = response.StatusCode
	Logger.Debug("Response StatusCode %d", r.statuscode)

	// Decode from r.response Accept() type
	if (success != nil || failure != nil) && r.statuscode != 204 {
		if 200 <= r.statuscode && r.statuscode <= 299 {
			if success != nil {

				if r.response != nil {
					Logger.Debug("Going to decode Response Body (%T) into success (%T)", r.response, success)

					if err := r.response.Decode(response.Body, success); err != nil {
						return r, err
					}
				}
			}
		} else {
			if failure != nil {

				if r.response != nil {
					Logger.Debug("Going to decode Response Body (%T) into failure (%T)", r.response, failure)

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

// BodyProvider sets the Request's body provider from original BodyProvider interface{}
func (r *Requist) BodyProvider(body BodyProvider) *Requist {

	Logger.Debug("Setting BodyProvider (%T)", body)

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

// BodyAsForm sets the Request's body from a formProvider
func (r *Requist) BodyAsForm(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(formProvider{payload: body})
}

// BodyAsJSON sets the Request's body from a jsonProvider
func (r *Requist) BodyAsJSON(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(jsonProvider{payload: body})
}

// BodyAsText sets the Request's body from a textProvider
func (r *Requist) BodyAsText(body interface{}) *Requist {

	if body == nil {
		return r
	}

	return r.BodyProvider(textProvider{payload: body})
}

//#$$=== Response Body functions, used to set type of response

// BodyResponse sets the response's body
func (r *Requist) BodyResponse(body BodyResponse) *Requist {

	Logger.Debug("Setting BodyResponse (%T)", body)

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

	Logger.Debug("Setting Accept (%s)", accept)

	switch accept {
	case FormContentType:
		r.BodyResponse(formResponse{})
	case JSONContentType:
		r.BodyResponse(jsonResponse{})
	case TextContentType:
		r.BodyResponse(textResponse{})
	default:
		r.response = nil
	}
}

//#$$=== QueryParams manipulation functions

// PrepareRequestURI returns actual uri
func (r *Requist) PrepareRequestURI() (uri string, err error) {

	if len(r.uri) > 0 {
		uri = r.uri
		r.uri = ""
		err = nil
	} else {
		reqURL, err := url.Parse(ParseBaseURL(r.url))
		if err != nil {
			return "", err
		}
		reqURL.Path = r.path
		reqURL.RawQuery = r.queries.Encode()
		uri = reqURL.String()
	}
	return uri, err
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

	r.url = ParseBaseURL(base)

	return r
}

// Path sets request path to use in next request
func (r *Requist) Path(path string) *Requist {

	r.path = ParsePathURL(r.url, path)

	return r
}

// URI sets request uri to use in next request
func (r *Requist) URI(uri string) *Requist {

	r.uri = uri

	return r
}

// Method set HTTP Method to execute
func (r *Requist) Method(method string) *Requist {

	switch strings.ToUpper(method) {
	case http.MethodGet:
		r.method = http.MethodGet

	case http.MethodPut:
		r.method = http.MethodPut

	case http.MethodPost:
		r.method = http.MethodPost

	case http.MethodPatch:
		r.method = http.MethodPatch

	case http.MethodDelete:
		r.method = http.MethodDelete

	case http.MethodOptions:
		r.method = http.MethodOptions

	case http.MethodConnect:
		r.method = http.MethodConnect

	default:
		r.method = http.MethodGet
	}

	return r
}

//#$$=== Requist functions executors, Correspond to HTTP Methods

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

// Connect implement CONNECT HTTP Method
func (r *Requist) Connect(path string, success, failure interface{}) (*Requist, error) {

	return r.Method(http.MethodConnect).Path(path).Request(success, failure)
}

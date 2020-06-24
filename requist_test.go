package requist

import (
	"encoding/base64"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// UserInfo, fictional user information
type UserInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UsersInfo, array of fictional user information
type UsersInfo struct {
	Users []UserInfo `json:"users"`
}

// GenericResponse
type GenericResponse struct {
	Result string `json:"result"`
}

func MockHTTPServer() *httptest.Server {
	// Mock http server
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/user" {
				switch r.Method {

				case http.MethodConnect:
					w.WriteHeader(http.StatusNoContent)

				case http.MethodGet:
					w.WriteHeader(http.StatusOK)
					w.Header().Add("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"users":[{"name": "Jonah Doe", "age": 50},{"name": "Jason Borne", "age": 47}]}`))

				case http.MethodPost:
					w.WriteHeader(http.StatusOK)
					w.Header().Add("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"result": "Created"}`))

				case http.MethodOptions:
					w.WriteHeader(http.StatusNoContent)
				}

			} else {
				if r.URL.Path == "/user/1000" {
					switch r.Method {

					case http.MethodGet:
						w.WriteHeader(http.StatusOK)
						w.Header().Add("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"name": "Jonah Doe", "age": 50}`))

					case http.MethodPatch:
						w.WriteHeader(http.StatusAccepted)
						w.Header().Add("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"result": "Accepted"}`))

					case http.MethodPut:
						w.WriteHeader(http.StatusAccepted)
						w.Header().Add("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"result": "Updated"}`))

					case http.MethodDelete:
						w.WriteHeader(http.StatusNoContent)
					}

				}
			}
		}),
	)
}

//===

func TestRequist_SetClientTransport(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to nil out default Transport
	emptyClient.SetClientTransport(nil)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if nil out default Transport
	assert.Nil(t, emptyClient.client.Transport)

	// We set again our default http Transport to cleanhttp.DefaultTransport()
	emptyClient.SetClientTransport(cleanhttp.DefaultTransport())

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if Not nil out default Transport
	assert.NotNil(t, emptyClient.client.Transport)
}

func TestRequist_SetClientTimeout(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var duration = 10 * time.Second

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to 10 seconds our default Timeout
	emptyClient.SetClientTimeout(duration)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if not empty our default Timeout
	assert.NotEmpty(t, emptyClient.client.Timeout)

	// We have a Timeout set?
	assert.Equal(t, duration, emptyClient.client.Timeout)
}

func TestRequist_BodyProvider(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set nil provider", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyProvider(nil)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, nil, emptyClient.provider)
	})

	t.Run("set JSON provider", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyProvider(jsonProvider{payload: nil})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a jsonProvider
		json := &jsonProvider{}

		// our data is correct?
		assert.EqualValues(t, JSONContentType, json.ContentType())
	})

	t.Run("set www-urlencoded provider", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyProvider(formProvider{payload: nil})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a formProvider
		form := &formProvider{}

		// our data is correct?
		assert.EqualValues(t, FormContentType, form.ContentType())
	})

	t.Run("set text/plain provider", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyProvider(formProvider{payload: nil})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a formProvider
		text := &textProvider{}

		// our data is correct?
		assert.EqualValues(t, TextContentType, text.ContentType())
	})
}

func TestRequist_BodyAsForm(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set nil Body", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyAsForm(nil)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, nil, emptyClient.provider)
	})

	t.Run("set not nil Body", func(t *testing.T) {

		userInfo := &UserInfo{Name: "Jonah Doe", Age: 47}

		// We set Accept Header to nothing
		emptyClient.BodyAsForm(userInfo)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Get the request Body
		body, _ := emptyClient.provider.Body()
		buffer := new(strings.Builder)
		_, err := io.Copy(buffer, body)

		// error getting Body?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, "Age=47&Name=Jonah+Doe", buffer.String())
	})
}

func TestRequist_BodyAsJSON(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set nil Body", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyAsJSON(nil)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, nil, emptyClient.provider)
	})

	t.Run("set not nil Body", func(t *testing.T) {

		userInfo := &UserInfo{Name: "Jonah Doe", Age: 47}

		// We set Accept Header to nothing
		emptyClient.BodyAsJSON(userInfo)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Get the request Body
		body, _ := emptyClient.provider.Body()
		buffer := new(strings.Builder)
		_, err := io.Copy(buffer, body)

		// error getting Body?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, "{\"name\":\"Jonah Doe\",\"age\":47}\n", buffer.String())
	})
}

func TestRequist_BodyAsText(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set nil Body", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyAsText(nil)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, nil, emptyClient.provider)
	})

	t.Run("set not nil Body", func(t *testing.T) {

		userInfo := &UserInfo{Name: "Jonah Doe", Age: 47}

		// We set Accept Header to nothing
		emptyClient.BodyAsText(userInfo)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Get the request Body
		body, _ := emptyClient.provider.Body()

		// our data is correct?
		assert.EqualValues(t, nil, body)
	})
}

func TestRequist_BodyResponse(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set nil response", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyResponse(nil)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, nil, emptyClient.provider)
	})

	t.Run("set JSON response", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyResponse(jsonResponse{})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a jsonProvider
		json := &jsonResponse{}

		// our data is correct?
		assert.EqualValues(t, JSONContentType, json.Accept())
	})

	t.Run("set www-urlencoded response", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyResponse(formResponse{})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a formProvider
		form := &formResponse{}

		// our data is correct?
		assert.EqualValues(t, FormContentType, form.Accept())
	})

	t.Run("set text/plain response", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.BodyProvider(formProvider{payload: nil})

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// Mock a formProvider
		text := &textResponse{}

		// our data is correct?
		assert.EqualValues(t, TextContentType, text.Accept())
	})
}

func TestRequist_PrepareRequestURI(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get empty URI if set empty Base and empty Path", func(t *testing.T) {

		emptyClient.Base("")
		emptyClient.Path("")

		result, err := emptyClient.PrepareRequestURI()
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, expected, result)
	})

	t.Run("get just base if set valid Base and empty Path", func(t *testing.T) {

		emptyClient.Base(baseURL)
		emptyClient.Path("")

		result, err := emptyClient.PrepareRequestURI()
		expected := baseURL

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, expected, result)
	})

	t.Run("get full URI if set valid Base and valid Path", func(t *testing.T) {

		// set the path to
		aPath := "/users"
		emptyClient.Base(baseURL)
		emptyClient.Path(aPath)

		result, err := emptyClient.PrepareRequestURI()
		expected := baseURL + aPath

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, expected, result)
	})

	t.Run("get complete URI if set valid Base and valid Path and QueryParams", func(t *testing.T) {

		// set the path to
		aPath := "/users"
		emptyClient.Base(baseURL)
		emptyClient.Path(aPath)
		emptyClient.AddQueryParam("name", "Jonah Doe")
		emptyClient.AddQueryParam("hobbies", "Bike")

		result, err := emptyClient.PrepareRequestURI()
		expected := baseURL + aPath + "?" + emptyClient.queries.Encode()

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.Nil(t, err)

		// our data is correct?
		assert.EqualValues(t, expected, result)
	})
}

//===

func TestRequist_Accept(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)

	t.Run("set empty Accept Header", func(t *testing.T) {
		// We set Accept Header to nothing
		emptyClient.Accept("")

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, "", emptyClient.header.Get(acceptHeader))
	})

	t.Run("set JSON Accept Header", func(t *testing.T) {
		// We set Accept Header to application/json
		emptyClient.Accept(JSONContentType)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, JSONContentType, emptyClient.header.Get(acceptHeader))
	})

	t.Run("set www-urlencoded Accept Header", func(t *testing.T) {
		// We set Accept Header to application/x-www-form-urlencoded
		emptyClient.Accept(FormContentType)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, FormContentType, emptyClient.header.Get(acceptHeader))
	})

	t.Run("set text/plain Accept Header", func(t *testing.T) {
		// We set Accept Header to text/plain
		emptyClient.Accept(TextContentType)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, TextContentType, emptyClient.header.Get(acceptHeader))
	})

}

func TestRequist_AddHeader(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var header = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to 10 seconds our default Timeout
	emptyClient.AddHeader(header["key"], header["value"])

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if not empty our default Timeout
	assert.NotNil(t, emptyClient.header)

	// if not empty our default Timeout
	assert.NotEmpty(t, emptyClient.header)

	// We have a Timeout set?
	assert.Equal(t, header["value"], emptyClient.header.Get(header["key"]))
}

func TestRequist_SetHeader(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var header = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to 10 seconds our default Timeout
	emptyClient.SetHeader(header["key"], header["value"])

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if not empty our default Timeout
	assert.NotNil(t, emptyClient.header)

	// if not empty our default Timeout
	assert.NotEmpty(t, emptyClient.header)

	// We have a Timeout set?
	assert.Equal(t, header["value"], emptyClient.header.Get(header["key"]))
}

func TestRequist_DelHeader(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var header = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set the Header", func(t *testing.T) {
		// We set one Header
		emptyClient.AddHeader(header["key"], header["value"])

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not empty our default Timeout
		assert.NotNil(t, emptyClient.header)

		// if not empty our default Timeout
		assert.NotEmpty(t, emptyClient.header)

		// We have a Timeout set?
		assert.Equal(t, header["value"], emptyClient.header.Get(header["key"]))
	})

	t.Run("delete the Header", func(t *testing.T) {
		// We delete the Header
		emptyClient.DelHeader(header["key"])

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not nil our Header repository
		assert.NotNil(t, emptyClient.header)

		// if empty our Header repository
		assert.Empty(t, emptyClient.header)

		// We have a empty Header with these Key
		assert.Equal(t, "", emptyClient.header.Get(header["key"]))
	})
}

func TestRequist_AddQueryParam(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var queryparams = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to 10 seconds our default Timeout
	emptyClient.AddQueryParam(queryparams["key"], queryparams["value"])

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if not empty our default Timeout
	assert.NotNil(t, emptyClient.queries)

	// if not empty our default Timeout
	assert.NotEmpty(t, emptyClient.queries)

	// We have a Timeout set?
	assert.Equal(t, queryparams["value"], emptyClient.queries.Get(queryparams["key"]))
}

func TestRequist_SetQueryParam(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var queryparams = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// We set to 10 seconds our default Timeout
	emptyClient.SetQueryParam(queryparams["key"], queryparams["value"])

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if not empty our default Timeout
	assert.NotNil(t, emptyClient.queries)

	// if not empty our default Timeout
	assert.NotEmpty(t, emptyClient.queries)

	// We have a Timeout set?
	assert.Equal(t, queryparams["value"], emptyClient.queries.Get(queryparams["key"]))
}

func TestRequist_DelQueryParam(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var queryparams = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set the QueryParam", func(t *testing.T) {
		// We set one Header
		emptyClient.AddQueryParam(queryparams["key"], queryparams["value"])

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not empty our default Timeout
		assert.NotNil(t, emptyClient.queries)

		// if not empty our default Timeout
		assert.NotEmpty(t, emptyClient.queries)

		// We have a Timeout set?
		assert.Equal(t, queryparams["value"], emptyClient.queries.Get(queryparams["key"]))
	})

	t.Run("delete the QueryParam", func(t *testing.T) {
		// We delete the Header
		emptyClient.DelQueryParam(queryparams["key"])

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not nil our Header repository
		assert.NotNil(t, emptyClient.queries)

		// if empty our Header repository
		assert.Empty(t, emptyClient.queries)

		// We have a empty Header with these Key
		assert.Equal(t, "", emptyClient.queries.Get(queryparams["key"]))
	})
}

func TestRequist_CleanQueryParams(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"
	var queryparams = map[string]string{"key": "X-Header", "value": "false"}

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("set the QueryParam", func(t *testing.T) {
		// We set one Header
		emptyClient.AddQueryParam(queryparams["key"], queryparams["value"])

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not empty our default Timeout
		assert.NotNil(t, emptyClient.queries)

		// if not empty our default Timeout
		assert.NotEmpty(t, emptyClient.queries)

		// We have a Timeout set?
		assert.Equal(t, queryparams["value"], emptyClient.queries.Get(queryparams["key"]))
	})

	t.Run("delete the QueryParam", func(t *testing.T) {
		// We delete the Header
		emptyClient.CleanQueryParams()

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if not nil our Header repository
		assert.NotNil(t, emptyClient.queries)

		// if empty our Header repository
		assert.Empty(t, emptyClient.queries)

		// We have a empty Header with these Key
		assert.Equal(t, "", emptyClient.queries.Get(queryparams["key"]))
	})
}

func TestRequist_SetBasicAuth(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get empty Auth if set empty Username and empty Password", func(t *testing.T) {

		// We set some variables
		username := ""
		password := ""

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.auth)
	})

	t.Run("get empty Auth if set valid Username and empty Password", func(t *testing.T) {

		// We set some variables
		username := "anonymous"
		password := ""

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.auth)
	})

	t.Run("get empty Auth if set empty Username and valid Password", func(t *testing.T) {

		// We set some variables
		username := ""
		password := "Password123"

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.auth)
	})

	t.Run("get valid Auth if set valid Username and valid Password", func(t *testing.T) {

		// We set some variables
		username := "anonymous"
		password := "Password123"

		emptyClient.SetBasicAuth(username, password)
		expectedPlain := "anonymous:Password123"
		expectedBase64 := "Basic " + base64.StdEncoding.EncodeToString([]byte(expectedPlain))

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expectedPlain, emptyClient.auth)
		assert.EqualValues(t, expectedBase64, emptyClient.header.Get("Authorization"))
	})
}

func TestRequist_StatusCode(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get 200 StatusCode", func(t *testing.T) {

		emptyClient.statuscode = 200
		expected := 200

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.StatusCode())
	})

	t.Run("get 504 StatusCode", func(t *testing.T) {

		emptyClient.statuscode = 504
		expected := 504

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.StatusCode())
	})
}

func TestRequist_GetBasicAuth(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get empty Auth if set empty Username and empty Password", func(t *testing.T) {

		// We set some variables
		username := ""
		password := ""

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.GetBasicAuth())
	})

	t.Run("get empty Auth if set valid Username and empty Password", func(t *testing.T) {

		// We set some variables
		username := "anonymous"
		password := ""

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.GetBasicAuth())
	})

	t.Run("get empty Auth if set empty Username and valid Password", func(t *testing.T) {

		// We set some variables
		username := ""
		password := "Password123"

		emptyClient.SetBasicAuth(username, password)
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.GetBasicAuth())
	})

	t.Run("get valid Auth if set valid Username and valid Password", func(t *testing.T) {

		// We set some variables
		username := "anonymous"
		password := "Password123"

		emptyClient.SetBasicAuth(username, password)
		expected := "anonymous:Password123"

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.GetBasicAuth())
	})
}

//===

func TestRequist_Base(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get empty Base if baseURL is empty", func(t *testing.T) {

		emptyClient.Base("")
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.url)
	})

	t.Run("get empty Base if baseURL is invalid", func(t *testing.T) {

		emptyClient.Base("file:///root/test/filename.json")
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.url)
	})

	t.Run("get same Base if baseURL is valid", func(t *testing.T) {

		emptyClient.Base(baseURL)
		expected := baseURL

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, baseURL)
	})
}

func TestRequist_Path(t *testing.T) {

	// We define some variables
	var baseURL = "http://live.apitest.org"

	// We create our requist Client
	emptyClient := New(baseURL)

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	t.Run("get empty Path if path is empty", func(t *testing.T) {

		emptyClient.Path("")
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.path)
	})

	t.Run("get empty Path if path is invalid", func(t *testing.T) {

		emptyClient.Path("file:///root/test/filename.json")
		expected := ""

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.path)
	})

	t.Run("get same Path if path is valid", func(t *testing.T) {

		path := "/valid/route/to/resource"
		emptyClient.Path(path)
		expected := path

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// our data is correct?
		assert.EqualValues(t, expected, emptyClient.path)
	})
}

func TestRequist_Method(t *testing.T) {}

//===

func TestRequist_New(t *testing.T) {

	t.Run("return nil if a empty baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = ""
		// fire up
		// We create our requist Client
		result := New(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Nil(t, result)
	})

	t.Run("return nil if a invalid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://?bar&?foo"
		// fire up
		// We create our requist Client
		result := New(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Nil(t, result)
	})

	t.Run("return new Client if a valid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "http://live.apitest.org"
		// fire up
		// We create our requist Client
		result := New(baseURL)
		// if result equals to expected?
		// we have a new Client?
		assert.NotNil(t, result)
	})

}

func TestRequist_Get(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	t.Run("return one resource via Get", func(t *testing.T) {
		// Define some vars
		success := &UserInfo{}
		var fail interface{}

		// We create our requist Client
		emptyClient := New(server.URL)

		// Set JSON response body
		emptyClient.Accept(JSONContentType)

		// fire up the request
		client, err := emptyClient.Get("/user/1000", success, fail)

		// We test the results
		// if client return not Nil?
		assert.Nil(t, err)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.NotNil(t, client)

		// are the same object?
		assert.EqualValues(t, emptyClient, client)

		// our data is correct?
		assert.EqualValues(t, "Jonah Doe", success.Name)
		assert.EqualValues(t, 50, success.Age)
	})

	t.Run("return > 1 resource via Get", func(t *testing.T) {
		// Define some vars
		success := &UsersInfo{}
		var fail interface{}

		// We create our requist Client
		emptyClient := New(server.URL)

		// Set JSON response body
		emptyClient.Accept(JSONContentType)

		// fire up the request
		client, err := emptyClient.Get("/user", success, fail)

		// We test the results
		// if client return not Nil?
		assert.Nil(t, err)

		// was modified out Client?
		assert.NotNil(t, emptyClient)

		// if client return not Nil?
		assert.NotNil(t, client)

		// are the same object?
		assert.EqualValues(t, emptyClient, client)

		// our data is correct?
		assert.EqualValues(t, "Jonah Doe", success.Users[0].Name)
		assert.EqualValues(t, 50, success.Users[0].Age)

		// our data is correct?
		assert.EqualValues(t, "Jason Borne", success.Users[1].Name)
		assert.EqualValues(t, 47, success.Users[1].Age)
	})
}

func TestRequist_Put(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	success := &GenericResponse{}
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(UserInfo{Name: "Jason Borne", Age: 40}).Put("/user/1000", success, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, "Updated", success.Result)
}

func TestRequist_Post(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	success := &GenericResponse{}
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(UserInfo{Name: "Jason Borne", Age: 40}).Post("/user", success, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, "Created", success.Result)
}

func TestRequist_Patch(t *testing.T) {

	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	success := &GenericResponse{}
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(UserInfo{Name: "Jason Doe"}).Patch("/user/1000", success, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, "Accepted", success.Result)
}

func TestRequist_Delete(t *testing.T) {
	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(UserInfo{Name: "Jason Borne", Age: 40}).Delete("/user/1000", nil, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, 204, client.StatusCode())
}

func TestRequist_Options(t *testing.T) {
	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(nil).Options("/user", nil, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, 204, client.StatusCode())
}

func TestRequist_Connect(t *testing.T) {
	// We create a Mock Server
	server := MockHTTPServer()
	defer server.Close()

	// We create our requist Client
	emptyClient := New(server.URL)
	emptyClient.Accept(JSONContentType)

	// Some variables
	var fail interface{}

	// fire up the request
	client, err := emptyClient.BodyAsJSON(nil).Connect("/user", nil, fail)

	// We test the results

	// was modified out Client?
	assert.NotNil(t, emptyClient)

	// if client return not Nil?
	assert.NotNil(t, client)

	// are the same object?
	assert.EqualValues(t, emptyClient, client)

	// if client return not Nil?
	assert.Nil(t, err)

	// our data is correct?
	assert.EqualValues(t, 204, client.StatusCode())
}

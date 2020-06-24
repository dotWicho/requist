package requist

import "time"

//=== Useful constants

const (
	acceptHeader string = "Accept"
	contentType  string = "Content-Type"
	// TextContentType is an alias to HTTP text/plain MIME Type
	TextContentType string = "text/plain"
	// JSONContentType is an alias to HTTP application/json MIME Type
	JSONContentType string = "application/json"
	// FormContentType is an alias to HTTP application/x-www-form-urlencoded MIME Type
	FormContentType string = "application/x-www-form-urlencoded"

	// Timeout of http.Client default, 4 seconds
	defaultTimeout = 4 * time.Second
)

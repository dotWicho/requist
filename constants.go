package requist

import "time"

//=== Useful constants

const (
	acceptHeader    string = "Accept"
	contentType     string = "Content-Type"
	TextContentType string = "text/plain"
	JsonContentType string = "application/json"
	FormContentType string = "application/x-www-form-urlencoded"

	// Timeout of http.Client default, 4 seconds
	defaultTimeout = 4 * time.Second
)

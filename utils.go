package requist

import (
	"log"
	"net/url"
)

//=== Supplemental functions to manipulate path

// ParseBaseURL check if is valid the base string passed
func parseBaseURL(base string) string {

	urlParsed, err := url.Parse(base)
	if err != nil {
		log.Fatalln()
	}
	urlParsed.RawQuery = ""
	urlParsed.Fragment = ""

	return urlParsed.String()
}

// ParsePathURL check relative path
func parsePathURL(base string, path string) string {

	urlParsed, err := url.Parse(base)
	if err != nil {
		log.Fatalln()
	}
	urlParsed.Path = path
	return urlParsed.String()
}

package requist

import (
	"net/url"
	"strings"
)

//=== Supplemental functions to manipulate path

// IsValidScheme check validity of a scheme
func IsValidScheme(scheme string) bool {

	return scheme != "" && scheme != "file" && (scheme == "http" || scheme == "https")
}

// IsValidHostname check validity of a hostname
func IsValidHostname(host string) bool {

	return host != "" && !strings.HasPrefix(host, ".") && !strings.HasSuffix(host, ":")
}

// IsValidBase check validity of a base URL
func IsValidBase(base string) bool {

	urlParsed, err := url.Parse(base)
	if err != nil || !IsValidScheme(urlParsed.Scheme) || !IsValidHostname(urlParsed.Host) {
		return false
	}

	return true
}

// ParseBaseURL check if is valid the base string passed
func ParseBaseURL(base string) string {

	urlParsed, err := url.Parse(base)
	if err != nil || !IsValidScheme(urlParsed.Scheme) || !IsValidHostname(urlParsed.Host) {
		return ""
	}
	urlParsed.RawQuery = ""
	urlParsed.Fragment = ""
	urlParsed.Path = ""
	urlParsed.Opaque = ""

	return urlParsed.String()
}

// ParsePathURL check relative path
func ParsePathURL(base string, path string) string {

	urlParsed, err := url.Parse(base + path)
	if err != nil || !IsValidScheme(urlParsed.Scheme) || !IsValidHostname(urlParsed.Host) {
		return ""
	}

	return urlParsed.Path
}

package requist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidScheme(t *testing.T) {

	t.Run("return false if an empty scheme", func(t *testing.T) {
		// Define some vars
		var scheme = ""

		// fire up
		// We create our requist Client
		result := IsValidScheme(scheme)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid scheme (file)", func(t *testing.T) {
		// Define some vars
		var scheme = "file"

		// fire up
		// We create our requist Client
		result := IsValidScheme(scheme)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return true if a valid scheme (http)", func(t *testing.T) {
		// Define some vars
		var scheme = "http"

		// fire up
		// We create our requist Client
		result := IsValidScheme(scheme)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, true, result)
	})

	t.Run("return true if a valid scheme (http)", func(t *testing.T) {
		// Define some vars
		var scheme = "https"

		// fire up
		// We create our requist Client
		result := IsValidScheme(scheme)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, true, result)
	})
}

func TestIsValidBase(t *testing.T) {

	t.Run("return false if an empty baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = ""

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://?bar&?foo"

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid scheme in baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "file:///home/user/config.txt"

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid host in baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://.ourhost/api/resource"

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid host in baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://ourhost:/api/resource"

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return true if a valid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://live.apitest.org"

		// fire up
		// We create our requist Client
		result := IsValidBase(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, true, result)
	})
}

func TestIsValidHostname(t *testing.T) {

	t.Run("return false if an empty hostname", func(t *testing.T) {
		// Define some vars
		var hostname = ""

		// fire up
		// We create our requist Client
		result := IsValidHostname(hostname)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid host in baseURL", func(t *testing.T) {
		// Define some vars
		var hostname = ".apitest.org"

		// fire up
		// We create our requist Client
		result := IsValidHostname(hostname)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return false if an invalid host in baseURL", func(t *testing.T) {
		// Define some vars
		var hostname = "live.apitest:"

		// fire up
		// We create our requist Client
		result := IsValidHostname(hostname)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, false, result)
	})

	t.Run("return true if a valid host in baseURL", func(t *testing.T) {
		// Define some vars
		var hostname = "ive.apitest.org"

		// fire up
		// We create our requist Client
		result := IsValidHostname(hostname)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, true, result)
	})
}

func TestParseBaseURL(t *testing.T) {

	t.Run("return empty string if a invalid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://?bar&?foo"
		var expected = ""

		// fire up
		// We create our requist Client
		result := ParseBaseURL(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, expected, result)
	})

	t.Run("return same if a valid baseURL", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://live.apitest.org/path/to/resource"
		var expected = "https://live.apitest.org"

		// fire up
		// We create our requist Client
		result := ParseBaseURL(baseURL)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, expected, result)
	})
}

func TestWithPathURL(t *testing.T) {

	t.Run("return empty string if a in valid baseURL and invalid path", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://https://https://google.com/"
		var path = "/¿test=+to/resource"
		var expected = ""

		// fire up
		// We create our requist Client
		result := ParsePathURL(baseURL, path)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, expected, result)
	})

	t.Run("return empty string if a in invalid baseURL and valid path", func(t *testing.T) {
		// Define some vars
		var baseURL = "file:///folder/mode"
		var path = "/test/resource"
		var expected = ""

		// fire up
		// We create our requist Client
		result := ParsePathURL(baseURL, path)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, expected, result)
	})

	t.Run("return fullUrl string if a valid baseURL and path", func(t *testing.T) {
		// Define some vars
		var baseURL = "https://live.apitest.org"
		var path = "/path/to/resource"
		var expected = path

		// fire up
		// We create our requist Client
		result := ParsePathURL(baseURL, path)
		// if result equals to expected?
		// we don't have a new Client?
		assert.Equal(t, expected, result)
	})
}

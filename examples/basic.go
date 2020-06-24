package main

import (
	"encoding/json"
	"fmt"
	"github.com/dotWicho/requist"
	"os"
)

// UserInfo, fictional user information
type UserInfo struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}

// GetResponse encapsulates the httpbin.org response to GET Method in JSON format
type GetResponse struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

// PostResponse encapsulates the httpbin.org response to POST Method in JSON format
type PostResponse struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Files   map[string]string `json:"files"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	JSON    interface{}       `json:"json"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

func asJSON(body interface{}) string {
	strJson, _ := json.MarshalIndent(body, "", "  ")
	return string(strJson)
}

func main() {
	// We create the client using the server base URL we want to access
	client := requist.New("https://httpbin.org")
	// We want the answer in JSON format
	client.Accept(requist.JSONContentType)

	// Instantiate where we will get the answer
	getSuccess := &GetResponse{}
	// One trick to encapsulate all responses ;-),
	// If you have a custom format used it here
	var fail interface{}

	// We launched the query, first set body to nil, and then fire Get method
	if _, err := client.BodyAsJSON(nil).Get("/get", getSuccess, fail); err != nil {
		// if there's an error, exit with a panic
		panic(err)
	}
	fmt.Printf("%s\n", asJSON(getSuccess))

	// Instantiate where we will get the answer
	postSuccess := &PostResponse{}

	// We launched the query, first set body to nil, and then fire Post method
	if _, err := client.BodyAsJSON(nil).Post("/post", postSuccess, fail); err != nil {
		// if there's an error, exit with a panic
		panic(err)
	}
	fmt.Printf("%s\n", asJSON(postSuccess))

	// Cleanups
	postSuccess = &PostResponse{}

	// Populate our User
	body := &UserInfo{
		Name:    "Jonah Doe",
		Age:     33,
		Hobbies: []string{"Bike", "Trekking", "Coding"},
	}

	// We launched the query, first set body to our user Jonah, and then fire Post method
	if _, err := client.BodyAsJSON(body).Post("/post", postSuccess, fail); err != nil {
		// if there's an error, exit with a panic
		panic(err)
	}
	fmt.Printf("%s\n", asJSON(postSuccess))

	os.Exit(0)
}

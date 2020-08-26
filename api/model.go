package main

// db object for storing endpoint
type Endpoint struct {
	ID  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

// http return object for endpoint generation
type EndpointOutput struct {
	AppUrl  string `json:"appurl,omitempty"`
	YourUrl string `json:"note,omitempty"`
}

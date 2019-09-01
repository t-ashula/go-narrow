package narrow

import (
	"fmt"
	"net/http"
)

// A Client fetch data from novel api
type Client struct {
	httpClient *http.Client
}

var userAgent = fmt.Sprintf("go-narrow/%s", Version)

// NewClient returns new novel api client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

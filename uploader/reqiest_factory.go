package uploader

import (
	"io"
	"net/http"
)

type RequestFactory func(method string, urlStr string, body io.Reader) (*http.Request, error)

type basicRequestFactory struct {
	username string
	password string
}

func NewBasicRequestFactory(username string, password string) RequestFactory {
	return func(method string, url string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(username, password)
		return req, nil
	}
}

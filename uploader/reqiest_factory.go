package uploader

import (
	"io"
	"net/http"
)

type RequestFactory interface {
	NewRequest(method string, urlStr string, body io.Reader) (*http.Request, error)
}

type basicRequestFactory struct {
	username string
	password string
}

func (this basicRequestFactory) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(this.username, this.password)
	return req, nil
}

func NewBasicRequestFactory(username string, password string) RequestFactory {
	return basicRequestFactory{username, password}
}

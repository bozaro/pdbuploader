package uploader

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Uploader struct {
	SupportPutOverwrite bool
	client              *http.Client
	reqFactory          RequestFactory
}

const (
	maxRetryCount int = 5
	initialRetryDelay int = 100
)

type requestProvider func() (*http.Request, error)

func NewUploader(client *http.Client, reqFactory RequestFactory) Uploader {
	return Uploader{
		SupportPutOverwrite: false,
		client:              client,
		reqFactory:          reqFactory,
	}
}

func parentUrl(u *url.URL) *url.URL {
	if u == nil || u.Path == "" || u.Path == "/" {
		return nil
	}
	r := *u
	r.Path = path.Dir(r.Path)
	return &r
}

func (this Uploader) doRequest(provider requestProvider) (*http.Response, error) {
	retryDelay := initialRetryDelay
	for pass := 0;; pass++ {
		request, err := provider()
		if err != nil {
			return nil, err
		}
		response, err := this.client.Do(request)
		if (err == nil) && (response.StatusCode < 500 || response.StatusCode >= 600) {
			return response, nil
		}
		if pass >= maxRetryCount {
			return response, err
		}
		time.Sleep(time.Duration(retryDelay) * time.Millisecond)
		retryDelay *= 2
	}
}

func (this Uploader) MkDir(u *url.URL) error {
	if u.Path == "" || u.Path == "/" {
		return errors.New(fmt.Sprintf("Can't create root path: %s", u.String()))
	}
	requestFunc := func() (*http.Request, error) {
		return this.reqFactory("MKCOL", u.String(), nil)
	}
	response, err := this.doRequest(requestFunc)
	if err != nil {
		return err
	}
	// Already created
	if response.StatusCode == 405 {
		return nil
	}
	if response.StatusCode == 409 {
		if err := this.MkDir(parentUrl(u)); err != nil {
			return err
		}
		response, err = this.doRequest(requestFunc)
		if err != nil {
			return err
		}
	}
	// Successfully created
	if response.StatusCode == 201 {
		return nil
	}
	return errors.New(fmt.Sprintf("Unexpected MKCOL status code %d [%s]: %s", response.StatusCode, response.Status, u.String()))
}

func (this Uploader) Exists(u *url.URL) (bool, error) {
	response, err := this.doRequest(func() (*http.Request, error) {
		return this.reqFactory("HEAD", u.String(), nil)
	})
	if err != nil {
		return false, err
	}
	switch response.StatusCode {
	// File already exists
	case 200:
		return true, nil
	// Not found
	case 404:
		return false, nil
	default:
		return false, errors.New(fmt.Sprintf("Unexpected HEAD status code %d [%s]: %s", response.StatusCode, response.Status, u.String()))
	}
}

func (this Uploader) UploadFile(u *url.URL, content ContentProvider, overwrite bool) error {
	if (u.Path == "") || strings.HasSuffix(u.Path, "/") {
		return errors.New(fmt.Sprintf("Invalid target URL. Require file name in path: %s", u.String()))
	}
	if !(overwrite && this.SupportPutOverwrite) {
		exists, err := this.Exists(u)
		if exists || err != nil {
			return err
		}
	}
	requestFunc := func() (*http.Request, error) {
		reader, err := content.GetReader()
		if err != nil {
			return nil, err
		}
		req, err := this.reqFactory("PUT", u.String(), reader)
		if err != nil {
			return nil, err
		}
		if !overwrite {
			req.Header.Set("Overwrite", "F")
		}
		return req, nil
	}
	response, err := this.doRequest(requestFunc)
	if err != nil {
		return err
	}
	if response.StatusCode == 409 {
		if err := this.MkDir(parentUrl(u)); err != nil {
			return err
		}
		response, err = this.doRequest(requestFunc)
		if err != nil {
			return err
		}
	}
	// Successfully uploaded
	if response.StatusCode == 201 {
		return nil
	}
	return errors.New(fmt.Sprintf("Unexpected PUT status code %d [%s]: %s", response.StatusCode, response.Status, u.String()))
}

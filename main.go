// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
// https://tech.yandex.ru/disk/doc/dg/reference/move-docpage/
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bozaro/pdbuploader/parse"
	"github.com/bozaro/pdbuploader/uploader"
	"github.com/bozaro/pdbuploader/wildcard"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const (
	maxRetryCount     int = 5
	initialRetryDelay int = 100
)

type requestProvider func() (*http.Request, error)

type Uploader struct {
	client     *http.Client
	reqFactory uploader.RequestFactory
}

func NewUploader(client *http.Client, reqFactory uploader.RequestFactory) Uploader {
	return Uploader{
		client,
		reqFactory,
	}
}

func (this Uploader) doRequest(provider requestProvider) (*http.Response, error) {
	retryDelay := initialRetryDelay
	for pass := 0; ; pass++ {
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

func parentUrl(u *url.URL) *url.URL {
	if u == nil || u.Path == "" || u.Path == "/" {
		return nil
	}
	r := *u
	r.Path = path.Dir(r.Path)
	return &r
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

func (this Uploader) UploadFile(u *url.URL, content uploader.ContentProvider) error {
	if strings.HasSuffix(u.Path, "/") {
		return errors.New(fmt.Sprintf("Invalid target URL. Require file name in path: %s", u.String()))
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
		req.Header.Set("Overwrite", "F") // todo
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

func upload(reqFactory uploader.RequestFactory) {
	http_uploader := NewUploader(http.DefaultClient, reqFactory)

	u, _ := url.Parse("https://webdav.yandex.ru/PDB/foo/bar/blah.txt")
	err := http_uploader.UploadFile(u, uploader.NewBytesContentProvider([]byte("Example")))
	fmt.Println(err)
}

func main() {

	wildcard.FindFiles(".", []string{
		"**.go",
		"!.*",
		"/tmp/test/*.go",
		"!foo.go",
	}, nil)

	usernamePtr := flag.String("username", "bozaro", "Username")
	passwordPtr := flag.String("password", "", "Password")
	flag.Parse()

	if *passwordPtr != "" {
		upload(uploader.NewBasicRequestFactory(*usernamePtr, *passwordPtr))
	}

	file, _ := os.Open("sample/hello.exe")
	{
		info, err := parse.ParseExe(file)
		fmt.Printf("EXE: %s\n", err)
		if err == nil {
			fmt.Printf("  Code ID: %s\n", info.CodeId)
			fmt.Printf("  Debug ID: %s\n", info.DebugId)
			fmt.Printf("  PDB: %s\n", info.PDBFileName)
		}
	}
	file, _ = os.Open("sample/hello.pdb")
	{
		debug_id, err := parse.ParsePdb(file)
		fmt.Printf("PDB: %s\n", err)
		if err == nil {
			fmt.Printf("  Debug ID: %s\n", debug_id)
		}
	}
}

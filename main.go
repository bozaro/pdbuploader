// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
// https://tech.yandex.ru/disk/doc/dg/reference/move-docpage/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bozaro/pdbuploader/parse"
	"github.com/bozaro/pdbuploader/uploader"
	"net/http"
	"os"
)

type Uploader  struct {
	client     *http.Client
	reqFactory uploader.RequestFactory
}

func NewUploader(client *http.Client, reqFactory uploader.RequestFactory) Uploader {
	return Uploader{
		client,
		reqFactory,
	}
}

func (this Uploader ) UploadFile(url string, content uploader.ContentProvider) error {
	for i := 0; i < 5; i++ {
		reader, err := content.GetReader()
		if err != nil {
			return err
		}
		request, err := this.reqFactory.NewRequest("PUT", url, reader)
		if err != nil {
			return err
		}
		response, err := this.client.Do(request)
		if err == nil {
			fmt.Println(response)
			break
		}
	}
	return nil
}

func upload(reqFactory uploader.RequestFactory) {
	http_uploader := NewUploader(http.DefaultClient, reqFactory)
	client := http.DefaultClient

	/*req, err := http.NewRequest("HEAD", "https://webdav.yandex.ru/PDB/test.txt", nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	fmt.Println(err)
	fmt.Println(resp)*/

	/*propfind := []byte("<?xml version=\"1.0\"?><propfind xmlns=\"DAV:\"><propname/></propfind>")
	req, err := http.NewRequest("PROPFIND", "https://webdav.yandex.ru/PDB", bytes.NewReader(propfind))
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(propfind)))
	req.Header.Set("Depth", "1")
	req.SetBasicAuth(username, password)*/
	{
		err := http_uploader.UploadFile("https://webdav.yandex.ru/PDB/foo/bar/blah.txt", uploader.NewBytesContentProvider([]byte("Example")))
		fmt.Println(err)
	}
	{
		req, _ := reqFactory.NewRequest("MKCOL", "https://webdav.yandex.ru/PDB/foo", nil)
		client.Do(req)
	}
	{
		data := []byte("Some file data")
		req, _ := reqFactory.NewRequest("PUT", "https://webdav.yandex.ru/PDB/foo/bar.txt~", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
		client.Do(req)
	}
	{
		req, err := reqFactory.NewRequest("MOVE", "https://webdav.yandex.ru/PDB/foo/bar.txt~", nil)
		req.Header.Set("Destination", "/PDB/foo/bar.txt")
		req.Header.Set("Overwrite", "T")

		resp, err := client.Do(req)
		fmt.Println(err)
		fmt.Println(resp)
	}

}

func main() {

	usernamePtr := flag.String("username", "bozaro", "Username")
	passwordPtr := flag.String("password", "", "Password")
	flag.Parse()

	upload(uploader.NewBasicRequestFactory(*usernamePtr, *passwordPtr))

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

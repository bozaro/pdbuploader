// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
// https://tech.yandex.ru/disk/doc/dg/reference/move-docpage/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bozaro/pdbuploader/parse"
	"net/http"
	"os"
)

func upload(username string, password string) {
	client := &http.Client{
	//    CheckRedirect: redirectPolicyFunc,
	}
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
		req, _ := http.NewRequest("MKCOL", "https://webdav.yandex.ru/PDB/foo", nil)
		req.SetBasicAuth(username, password)
		client.Do(req)
	}
	{
		data := []byte("Some file data")
		req, _ := http.NewRequest("PUT", "https://webdav.yandex.ru/PDB/foo/bar.txt~", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
		req.SetBasicAuth(username, password)
		client.Do(req)
	}
	{
		req, err := http.NewRequest("MOVE", "https://webdav.yandex.ru/PDB/foo/bar.txt~", nil)
		req.Header.Set("Destination", "/PDB/foo/bar.txt")
		req.Header.Set("Overwrite", "T")
		req.SetBasicAuth(username, password)

		resp, err := client.Do(req)
		fmt.Println(err)
		fmt.Println(resp)
	}

}

func main() {

	usernamePtr := flag.String("username", "bozaro", "Username")
	passwordPtr := flag.String("password", "", "Password")
	flag.Parse()

	upload(*usernamePtr, *passwordPtr)

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

// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
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

	propfind := []byte("<?xml version=\"1.0\"?><propfind xmlns=\"DAV:\"><propname/></propfind>")
	req, err := http.NewRequest("PROPFIND", "https://webdav.yandex.ru/PDB", bytes.NewReader(propfind))
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(propfind)))
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	fmt.Println(err)
	fmt.Println(resp)

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

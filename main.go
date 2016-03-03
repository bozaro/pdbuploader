// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
// https://tech.yandex.ru/disk/doc/dg/reference/move-docpage/
package main

import (
	"flag"
	"fmt"
	"github.com/bozaro/pdbuploader/parse"
	"github.com/bozaro/pdbuploader/uploader"
	"net/http"
	"net/url"
	"os"
)

func upload(reqFactory uploader.RequestFactory) {
	http_uploader := uploader.NewUploader(http.DefaultClient, reqFactory)

	u, _ := url.Parse("https://webdav.yandex.ru/PDB/foo/bar/blah.txt")
	err := http_uploader.UploadFile(u, uploader.NewBytesContentProvider([]byte("Example")), false)
	fmt.Println(err)
}

func main() {
	usernamePtr := flag.String("username", "bozaro", "Username")
	passwordPtr := flag.String("password", "", "Password")
	flag.Parse()

	upload(uploader.NewBasicRequestFactory(*usernamePtr, *passwordPtr))

	{
		file, _ := os.Open("sample/hello.exe")
		defer file.Close()

		info, err := parse.ParseExe(file)
		fmt.Printf("EXE: %s\n", err)
		if err == nil {
			fmt.Printf("  Code ID: %s\n", info.CodeId)
			fmt.Printf("  Debug ID: %s\n", info.DebugId)
			fmt.Printf("  PDB: %s\n", info.PDBFileName)
		}
	}
	{
		file, _ := os.Open("sample/hello.pdb")
		defer file.Close()

		debug_id, err := parse.ParsePdb(file)
		fmt.Printf("PDB: %s\n", err)
		if err == nil {
			fmt.Printf("  Debug ID: %s\n", debug_id)
		}
	}
}

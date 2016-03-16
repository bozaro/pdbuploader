package wildcard

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FindFiles(base_dir string, wildcards []string, only_from_base bool) {
	if !filepath.IsAbs(base_dir) {
		current_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		base_dir = path.Join(current_dir, base_dir)
	}
	base_dir = path.Clean(base_dir)
	if base_dir == "." {
		base_dir = "/"
	}

	for _, wildcard := range wildcards {
		negative := false
		if strings.HasPrefix(wildcard, "!") {
			wildcard = wildcard[1:]
			negative = true
		}
		if !filepath.IsAbs(wildcard) {
			wildcard = path.Join(base_dir, wildcard)
		}
		fmt.Println(negative, wildcard)
	}

	fmt.Println(base_dir)
}

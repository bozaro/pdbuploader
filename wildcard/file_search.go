package wildcard

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type find_state struct {
}

func FindFiles(base_dir string, wildcards []string, work_dir *string) error {
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

	if work_dir != nil {
		if !filepath.IsAbs(*work_dir) {
			current_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}
			dir := path.Join(current_dir, *work_dir)
			work_dir = &dir
		}
	}

	volumes := make(map[string]struct{}, len(wildcards))
	state := find_state{}
	for _, wildcard := range wildcards {
		negative := false
		if strings.HasPrefix(wildcard, "!") {
			wildcard = wildcard[1:]
			negative = true
		}
		if !filepath.IsAbs(wildcard) {
			wildcard = path.Join(base_dir, wildcard)
		}
		matcher, err := NewPathMatcher(wildcard)
		if err != nil {
			return err
		}
		volume := filepath.VolumeName(wildcard)
		volumes[volume] = struct{}{}
		fmt.Println(negative, wildcard, matcher)
	}

	if work_dir != nil {
		find_recursive(*work_dir, state)
	} else {
		for key := range volumes {
			find_recursive(key, state)
		}
	}

	fmt.Println(volumes)
	fmt.Println(base_dir)
	fmt.Println(work_dir)
	return nil
}

func find_recursive(work_dir string, state find_state) {

}

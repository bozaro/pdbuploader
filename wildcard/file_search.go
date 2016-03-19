package wildcard

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type FindWildcard struct {
	negative bool
	matcher  PathMatcher
}

type FindState struct {
	wildcards []FindWildcard
}

type FileInfoByName []os.FileInfo

func (this FileInfoByName) Len() int {
	return len(this)
}
func (this FileInfoByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this FileInfoByName) Less(i, j int) bool {
	return this[i].Name() < this[j].Name()
}

func (this FindState) CreateChild(name string, dir bool) *FindState {
	if len(this.wildcards) == 0 {
		return nil
	}
	wildcards := make([]FindWildcard, 0)

	last_negate := false
	for _, wildcard := range this.wildcards {
		matcher := wildcard.matcher.CreateChild(name, dir)
		if matcher != nil {
			last_negate = wildcard.negative && matcher.Matched()
			wildcards = append(wildcards, FindWildcard{wildcard.negative, matcher})
		}
	}

	if last_negate {
		return nil
	}
	if len(wildcards) == 0 {
		return nil
	}
	return &FindState{wildcards}
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
	state := FindState{make([]FindWildcard, 0)}
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
		state.wildcards = append(state.wildcards, FindWildcard{negative, matcher})
		fmt.Println(negative, wildcard, matcher)
	}

	if work_dir != nil {
		find_work_dir(*work_dir, state)
	} else {
		for key := range volumes {
			find_work_dir(key, state)
		}
	}

	fmt.Println(volumes)
	fmt.Println(base_dir)
	fmt.Println(work_dir)
	return nil
}

func find_work_dir(work_dir string, state FindState) {
	fmt.Println("+++++++++++++++")
	fmt.Println("WORK: ", work_dir)
	var cur_state *FindState = &state
	for _, item := range filepath.SplitList(work_dir) {
		cur_state = cur_state.CreateChild(item, true)
		if cur_state == nil {
			break
		}
		fmt.Println(item)
	}
	fmt.Println("---------------")
	fmt.Println(cur_state)
	find_recursive(work_dir+string(filepath.Separator), cur_state)
}

func find_recursive(work_dir string, state *FindState) {
	if state == nil {
		return
	}
	dirs := make(map[string]struct{}, len(state.wildcards))
	fast := true
	for _, wildcard := range state.wildcards {
		if !wildcard.negative {
			names := wildcard.matcher.ExactNames()
			if names == nil {
				fast = false
				break
			}
			for _, name := range *names {
				dirs[name] = struct{}{}
			}
		}
	}

	if fast {
		fmt.Println("FAST", work_dir)
		for key := range dirs {
			find_recursive(filepath.Join(work_dir, key), state.CreateChild(key, true))
		}
	} else {
		fmt.Println("SCAN", work_dir)
		if files, err := ioutil.ReadDir(work_dir); err == nil {
			sort.Sort(FileInfoByName(files))
			for _, file := range files {
				child := state.CreateChild(file.Name(), file.IsDir())
				if child != nil {
					if file.IsDir() {
						find_recursive(filepath.Join(work_dir, file.Name()), child)
					} else {
						fmt.Println("===>", filepath.Join(work_dir, file.Name()))
					}
				}
			}
		}
	}

	//fmt.Println(work_dir, state)
}

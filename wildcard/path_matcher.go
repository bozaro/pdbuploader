package wildcard

type PathMatcher interface {
	CreateChild(name string, dir bool) PathMatcher
	Matched() bool
}

func NewPathMatcher(pattern string) (PathMatcher, error) {
	nameMatchers, err := createNameMatchers(pattern)
	if err != nil {
		return nil, err
	}
	if len(nameMatchers) == 0 {
		return AlwaysMatcher{}, nil
	}
	if hasRecursive(nameMatchers) {
		if len(nameMatchers) == 2 && nameMatchers[0].Recursive() {
			return newFileMaskMatcher(nameMatchers[1]), nil
		} else {
			return newRecursivePathMatcher(nameMatchers), nil
		}
	} else {
		return newSimplePathMatcher(nameMatchers), nil
	}
}

func hasRecursive(nameMatchers []NameMatcher) bool {
	for _, matcher := range nameMatchers {
		if matcher.Recursive() {
			return true
		}
	}
	return false
}

func createNameMatchers(pattern string) ([]NameMatcher, error) {
	tokens := NormalizePattern(SplitPattern(pattern))
	result := make([]NameMatcher, len(tokens)-1)
	for i := range result {
		var err error
		result[i], err = NewNameMatcher(tokens[i+1])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type AlwaysMatcher struct {
}

func (this AlwaysMatcher) CreateChild(name string, dir bool) PathMatcher {
	return &AlwaysMatcher{}
}

func (this AlwaysMatcher) Matched() bool {
	return true
}

// Complex full-feature pattern matcher.
type FileMaskMatcher struct {
	matcher NameMatcher
}

func newFileMaskMatcher(matcher NameMatcher) FileMaskMatcher {
	return FileMaskMatcher{
		matcher: matcher,
	}
}

func (this FileMaskMatcher) CreateChild(name string, dir bool) PathMatcher {
	if this.matcher.Matched(name, dir) {
		return &AlwaysMatcher{}
	}
	if !dir {
		return nil
	}
	return this
}

func (this FileMaskMatcher) Matched() bool {
	return false
}

// Complex full-feature pattern matcher.
type RecursivePathMatcher struct {
	indexes      []int
	nameMatchers []NameMatcher
}

func newRecursivePathMatcher(nameMatchers []NameMatcher) RecursivePathMatcher {
	return RecursivePathMatcher{
		indexes:      []int{0},
		nameMatchers: nameMatchers,
	}
}

func (this RecursivePathMatcher) CreateChild(name string, dir bool) PathMatcher {
	childs := make([]int, len(this.indexes)*2)
	changed := false
	count := 0
	for _, index := range this.indexes {
		if this.nameMatchers[index].Matched(name, dir) {
			if this.nameMatchers[index].Recursive() {
				childs[count] = index
				count++
				if this.nameMatchers[index+1].Matched(name, dir) {
					if index+2 == len(this.nameMatchers) {
						return AlwaysMatcher{}
					}
					childs[count] = index + 2
					count++
					changed = true
				}
			} else {
				if index+1 == len(this.nameMatchers) {
					return AlwaysMatcher{}
				}
				childs[count] = index + 1
				count++
				changed = true
			}
		} else {
			changed = true
		}
	}
	if !dir {
		return nil
	}
	if !changed {
		return this
	}
	if count > 0 {
		return RecursivePathMatcher{
			nameMatchers: this.nameMatchers,
			indexes:      childs[:count],
		}
	} else {
		return nil
	}
}

func (this RecursivePathMatcher) Matched() bool {
	return false
}

// Matcher for patterns without "**".
type SimplePathMatcher struct {
	index        int
	nameMatchers []NameMatcher
}

func newSimplePathMatcher(nameMatchers []NameMatcher) SimplePathMatcher {
	return SimplePathMatcher{
		nameMatchers: nameMatchers,
		index:        0,
	}
}

func (this SimplePathMatcher) CreateChild(name string, dir bool) PathMatcher {
	if this.nameMatchers[this.index].Matched(name, dir) {
		if this.index+1 == len(this.nameMatchers) {
			return AlwaysMatcher{}
		}
		return SimplePathMatcher{
			nameMatchers: this.nameMatchers,
			index:        this.index + 1,
		}
	}
	return nil
}

func (this SimplePathMatcher) Matched() bool {
	return false
}

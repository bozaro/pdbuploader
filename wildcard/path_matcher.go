package wildcard

type PathMatcher interface {
	CreateChild(name string, dir bool) PathMatcher
	Matched() bool
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

func (this FileMaskMatcher) CreateChild(name string, dir bool) PathMatcher {
	if this.matcher.Matched(name, dir) {
		return &AlwaysMatcher{}
	}
	return this
}

func (this FileMaskMatcher) Matched() bool {
	return true
}

// Complex full-feature pattern matcher.
type RecursivePathMatcher struct {
	indexes      []int
	nameMatchers []NameMatcher
	selfMatch    bool
}

func newRecursivePathMatcher(nameMatchers []NameMatcher) RecursivePathMatcher {
	return RecursivePathMatcher{
		indexes:      []int{0},
		nameMatchers: nameMatchers,
		selfMatch:    recursiveMatched(nameMatchers, []int{0}),
	}
}

func (this RecursivePathMatcher) CreateChild(name string, dir bool) PathMatcher {
	childs := make([]int, len(this.indexes)*2)
	changed := false
	childMatch := false
	count := 0
	for _, index := range this.indexes {
		if this.nameMatchers[index].Matched(name, dir) {
			if this.nameMatchers[index].Recursive() {
				childs[count] = index
				count++
				if index+1 < len(this.nameMatchers) && this.nameMatchers[index+1].Matched(name, dir) {
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
				} else if index+2 == len(this.nameMatchers) && this.nameMatchers[index+1].Recursive() {
					childMatch = true
				}
				childs[count] = index + 1
				count++
				changed = true
			}
		} else {
			changed = true
		}
	}
	if !changed {
		return this
	}
	if count > 0 {
		return RecursivePathMatcher{
			nameMatchers: this.nameMatchers,
			indexes:      childs[:count],
			selfMatch:    childMatch,
		}
	} else {
		return nil
	}
}

func (this RecursivePathMatcher) Matched() bool {
	return this.selfMatch
}

func recursiveMatched(nameMatchers []NameMatcher, indexes []int) bool {
	if len(nameMatchers) > 0 && nameMatchers[len(nameMatchers)-1].Recursive() {
		for _, index := range indexes {
			if index == len(nameMatchers)-1 {
				return true
			}
		}
	}
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
		return SimplePathMatcher{nameMatchers: this.nameMatchers,
			index: this.index + 1}
	}
	return nil
}

func (this SimplePathMatcher) Matched() bool {
	return true
}

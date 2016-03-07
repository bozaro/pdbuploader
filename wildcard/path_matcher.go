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

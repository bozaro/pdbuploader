package wildcard

const PATH_SEPARATOR rune = '/'

/**
  Split pattern with saving slashes.

  @param pattern Path pattern.
  @return Path pattern items.
*/
func SplitPattern(path string) []string {
	result := []string{}
	start := 0
	for i, c := range path {
		if c == PATH_SEPARATOR {
			result = append(result, path[start:i+1])
			start = i + 1
		}
	}
	if start != len(path) {
		result = append(result, path[start:])
	}
	return result
}

/**
  Remove redundant pattern parts and make patterns more simple.

  @param tokens Original modifiable list.
  @return Return tokens,
*/
func NormalizePattern(tokens []string) []string {
	if len(tokens) == 1 {
		if tokens[0] == "/" {
			tokens[0] = "**/"
		} else {
			tokens = []string{"**/", tokens[0]}
		}
	}
	if len(tokens) == 0 || tokens[0] != "/" {
		tokens = append([]string{"/"}, tokens...)
	}
	return tokens
}

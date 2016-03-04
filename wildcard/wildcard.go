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

func NormalizePattern(tokens []string) []string {
	return tokens
}

package wildcard

import (
	"strings"
)

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
	// Copy array
	tokens = append([]string(nil), tokens...)
	// By default without slashes using mask for files in all subdirectories
	if len(tokens) == 1 && !strings.Contains(tokens[0], "/") {
		tokens = []string{"**/", tokens[0]}
	}
	// Normalized pattern always starts with "/"
	if len(tokens) == 0 || tokens[0] != "/" {
		tokens = append([]string{"/"}, tokens...)
	}
	// Replace:
	//  * "**/*/" to "*/**/"
	//  * "**/**/" to "**/"
	//  * "**.foo" to "**/*.foo"
	for i := 1; i < len(tokens); {
		thisToken := tokens[i]
		prevToken := tokens[i-1]
		if thisToken == "/" {
			tokens = append(tokens[:i], tokens[i+1:]...)
			continue
		}
		if thisToken == "**/" && prevToken == "**/" {
			tokens = append(tokens[:i], tokens[i+1:]...)
			continue
		}
		if thisToken != "**/" && strings.HasPrefix(thisToken, "**") {
			tokens = append(tokens[:i], append([]string{"**/", thisToken[1:]}, tokens[i+1:]...)...)
			continue
		}
		if thisToken == "*/" && prevToken == "**/" {
			tokens[i-1] = "*/"
			tokens[i] = "**/"
			i--
			continue
		}
		i++
	}
	// Remove tailing "**/" and "*"
	for len(tokens) > 0 {
		token := tokens[len(tokens)-1]
		if token == "**/" || token == "*" {
			tokens = tokens[:len(tokens)-1]
		} else {
			break
		}
	}
	return tokens
}

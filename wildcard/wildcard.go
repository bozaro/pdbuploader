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
	// Use "**.foo" as "**/*.foo"
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token != "**/" && strings.HasPrefix(token, "**") {
			tokens = append(tokens[:i], append([]string{"**/", token[1:]}, tokens[i+1:]...)...)
		}
	}
	// Replace:
	//  * "**/*/" to "*/**/"
	//  * "**/**/" to "**/"
	for i := 0; i < len(tokens)-1; {
		if i > 0 && tokens[i] == "/" {
			tokens = append(tokens[:i], tokens[i+1:]...)
		} else if tokens[i] == "**/" && tokens[i+1] == "**/" {
			tokens = append(tokens[:i], tokens[i+1:]...)
		} else if tokens[i] == "**/" && tokens[i+1] == "*/" {
			tokens[i] = "*/"
			tokens[i+1] = "**/"
			if i > 0 {
				i = i - 1
			}
		} else {
			i++
		}
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

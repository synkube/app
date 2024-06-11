package common

import "strings"

func BuildInfo() string {
	return "Build info: Version 1.0.0"
}

func AddIndent(input, indent string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

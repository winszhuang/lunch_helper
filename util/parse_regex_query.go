package util

import (
	"fmt"
	"regexp"
)

func ParseRegexQuery(input string, re *regexp.Regexp) []string {
	fmt.Println("---------------")
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return []string{}
	}

	result := []string{}
	for i, m := range matches {
		fmt.Println(m)
		if i == 0 {
			continue
		}
		result = append(result, m)
	}

	return result
}

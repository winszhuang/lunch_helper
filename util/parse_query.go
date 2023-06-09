package util

import (
	"fmt"
	"net/url"
	"regexp"
)

func ParseQuery(input string) ([]string, error) {
	values, err := url.ParseQuery(input)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for key, value := range values {
		if len(value) < 1 {
			return nil, fmt.Errorf("no value for key %s", key)
		}
		result = append(result, value[0])
	}
	return result, nil
}

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

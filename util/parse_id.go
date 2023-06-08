package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseId(key string, data string) (int, error) {
	// e.g. /restaurantmenu=、/food=
	str := fmt.Sprintf("/%s=", key)
	chunk := strings.Split(data, str)
	if len(chunk) != 2 {
		return 0, fmt.Errorf("參數帶入有問題")
	}
	return strconv.Atoi(chunk[1])
}

package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const charset = "0123456789abcdefghijklmnopqrstuvwxyz"

func init() {
	// 種子設置為當前時間的 UnixNano 格式
	// 為了保證每次運行程序時，隨機生成的數字序列都是不同的
	rand.Seed(time.Now().UnixNano())
}

func RandomLineID() string {
	rand.Seed(time.Now().UnixNano())

	id := make([]byte, 50)
	for i := 0; i < 50; i++ {
		id[i] = charset[rand.Intn(len(charset))]
	}

	return "C" + string(id)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomChar(n int) string {
	var sb strings.Builder
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

func RandomPicture() string {
	return fmt.Sprintf("https://loremflickr.com/320/240/%s", RandomChar(40))
}

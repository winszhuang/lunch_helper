package util

import "unicode/utf8"

// 確保切割後的字串不會出現亂碼
func TruncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) > maxLen {
		rs := []rune(s)
		if len(rs) > maxLen {
			return string(rs[:maxLen]) + "..."
		}
	}
	return s
}

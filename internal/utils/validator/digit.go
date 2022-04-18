package validator

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

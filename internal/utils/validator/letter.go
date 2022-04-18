package validator

func IsUpperLatin(ch rune) bool {
	return 'A' <= ch && ch <= 'Z'
}

func IsLowerLatin(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}

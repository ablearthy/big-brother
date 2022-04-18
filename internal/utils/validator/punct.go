package validator

func IsPunct(ch rune) bool {
	return ('!' <= ch && ch <= '/') || (':' <= ch && ch <= '@') || ('[' <= ch && ch <= '`') || ('{' <= ch && ch <= '~')
}

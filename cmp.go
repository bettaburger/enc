package enc

func StrCmp(a string, b string) int {
	switch {
	// if equal return 0
	case a == b:
		return 0
		// otherwise return -1
	default:
		return -1
	}
}

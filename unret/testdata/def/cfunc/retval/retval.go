package cfretval

func rets3() (int, string, error) { return 4, "5", nil }

func use() (string, int, error) {
	a, b, c := rets3()
	return b, a, c
}

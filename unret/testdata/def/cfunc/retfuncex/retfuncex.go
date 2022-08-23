package cfretfuncex

func rets3() (int, string, error) { return 4, "5", nil }

func Use() (int, string, error) {
	return rets3()
}

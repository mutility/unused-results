package cfretvaltup

func rets3() (int, string, error) { return 4, "5", nil }

func use() (int, string, error) {
	return rets3()
}

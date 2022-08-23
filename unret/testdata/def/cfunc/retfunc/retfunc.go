package cfretfunc

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }
func rets3() (int, string, error) { return 4, "5", nil }

func use() (func(), func() int, func() (int, int), func() (int, string, error)) {
	return rets0, rets1, rets2, rets3
}

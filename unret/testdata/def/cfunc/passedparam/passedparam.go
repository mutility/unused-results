package cfpassedparam

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }
func rets3() (int, string, error) { return 4, "5", nil }

func use(fn func(...interface{})) {
	fn(rets0)
	fn(rets1)
	fn(rets2)
	fn(rets3)
}

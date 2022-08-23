package cspassedparam

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use(s *struc, fn func(...interface{})) {
	fn(s.rets0)
	// TODO: fn(s.rets1)
	// TODO: fn(s.rets2)
	// TODO: fn(s.rets3)
}

package cspassedparam

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func Use(s *struc, fn func(...interface{})) {
	fn(s.rets0)
	fn(s.rets1)
	fn(s.rets2)
	fn(s.rets3)
}

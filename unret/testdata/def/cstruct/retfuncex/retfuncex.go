package csretfuncex

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

/* TODO:
func Use() (func(), func() int, func() (int, int), func() (int, string, error)) {
	var s struc
	return s.rets0, s.rets1, s.rets2, s.rets3
}
*/

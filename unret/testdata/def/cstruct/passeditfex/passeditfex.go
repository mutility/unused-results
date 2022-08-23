package cspasseditfex

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use(s *struc, itf interface{ Use(...interface{}) }) {
	itf.Use(s.rets0)
	// TODO: itf.Use(s.rets1)
	// TODO: itf.Use(s.rets2)
	// TODO: itf.Use(s.rets3)
}

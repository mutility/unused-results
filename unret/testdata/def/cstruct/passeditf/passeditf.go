package cspasseditf

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use(s *struc, itf interface{ use(...interface{}) }) {
	itf.use(s.rets0)
	itf.use(s.rets1)
	itf.use(s.rets2)
	itf.use(s.rets3)
}

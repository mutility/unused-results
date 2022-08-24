package ppasseditf

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }           // want "0 .int. is never used"
func (*struc) rets2() (x, y int)           { return 2, 3 }        // want "0 .x int. is never used" "1 .y int. is never used"
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // want "0 .int. is never used" "1 .string. is never used" "2 .error. is never used"

func use(s *struc, itf interface{ use(...interface{}) }) {
	itf.use(s.rets0)
	itf.use(s.rets1)
	itf.use(s.rets2)
	itf.use(s.rets3)
}

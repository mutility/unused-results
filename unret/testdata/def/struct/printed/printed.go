package sprinted

import "log"

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func (*struc) rets2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use() {
	s := struc{}

	s.rets1()
	a, _ := s.rets2()
	_, b, _ := s.rets3()
	_, _, c := s.rets3()

	log.Print(a, b, c)
}

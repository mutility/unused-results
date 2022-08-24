package sunexported

import "log"

type struc struct {
	a, b, c int
}

//  since struc isn't exported and isn't returned, even its exported methods shouldn't be marked used

func (s *struc) Res0()                      {}
func (s *struc) Res1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func (s *struc) Res2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func (s *struc) Res3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use(s *struc) {
	s.Res0()
	s.Res1()

	a, _ := s.Res2()
	_, b, _ := s.Res3()
	_, _, c := s.Res3()

	log.Print(a, b, c)
}

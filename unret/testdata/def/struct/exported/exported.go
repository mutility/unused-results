package sexported

import "log"

type S struct {
	a, b, c int
}

func (s *S) res0()                      {}
func (s *S) res1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func (s *S) res2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func (s *S) res3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use(s *S) {
	s.res0()
	s.res1()

	a, _ := s.res2()
	_, b, _ := s.res3()
	_, _, c := s.res3()

	log.Print(a, b, c)
}

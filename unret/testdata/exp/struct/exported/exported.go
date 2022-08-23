package xsexported

import "log"

type S struct {
	a, b, c int
}

func (s *S) res0()
func (s *S) res1() int                  // want "result 0 \\(int\\) is never used"
func (s *S) res2() (x, y int)           // want "result 1 \\(y int\\) is never used"
func (s *S) res3() (int, string, error) // want "result 0 \\(int\\) is never used"
func (s *S) Res0()
func (s *S) Res1() int                  // want "result 0 \\(int\\) is never used"
func (s *S) Res2() (x, y int)           // want "result 1 \\(y int\\) is never used"
func (s *S) Res3() (int, string, error) // want "result 0 \\(int\\) is never used"

func use(s *S) {
	s.Res0()
	s.res0()
	s.Res1()
	s.res1()

	a, _ := s.res2()
	_, b, _ := s.res3()
	_, _, c := s.res3()

	log.Print(a, b, c)

	a, _ = s.Res2()
	_, b, _ = s.Res3()
	_, _, c = s.Res3()

	log.Print(a, b, c)
}

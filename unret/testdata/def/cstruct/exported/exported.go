package csexported

import "log"

type S struct {
	a, b, c int
}

func (s *S) Res0()
func (s *S) Res1() int
func (s *S) Res2() (x, y int)
func (s *S) Res3() (int, string, error)

func use(s *S) {
	s.Res0()
	s.Res1()

	a, _ := s.Res2()
	_, b, _ := s.Res3()
	_, _, c := s.Res3()

	log.Print(a, b, c)
}

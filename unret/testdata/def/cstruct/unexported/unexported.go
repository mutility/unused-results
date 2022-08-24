package sunexported

import "log"

type struc struct {
	a, b, c int
}

// since *struc is returned by an exported function, assume its exported methods are used

func (s *struc) Res0()                      {}
func (s *struc) Res1() int                  { return 1 }
func (s *struc) Res2() (x, y int)           { return 2, 3 }
func (s *struc) Res3() (int, string, error) { return 4, "5", nil }

func Use(s *struc) *struc {
	s.Res0()
	s.Res1()

	a, _ := s.Res2()
	_, b, _ := s.Res3()
	_, _, c := s.Res3()

	log.Print(a, b, c)
	return s
}

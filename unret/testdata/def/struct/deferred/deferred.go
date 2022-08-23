package sdeferred

type struc struct{ a, b, c int }

func (s *struc) vals() (a, b, c int) { return s.a, s.b, s.c } // want "1 \\(b int\\) is never used"

func use(s *struc) {
	x, _, _ := s.vals()
	defer func(int) {
		x, _, z := s.vals()
		println(x + z)
	}(x)
}

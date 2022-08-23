package csdeferred

type struc struct{ a, b, c int }

func (s *struc) vals() (a, b, c int) { return s.a, s.b, s.c }

func use(s *struc) {
	x, _, _ := s.vals()
	defer func(int) {
		_, y, z := s.vals()
		println(y + z)
	}(x)
}

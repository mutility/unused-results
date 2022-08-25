package cfstruct

type struc struct{ x, y int }

func (s *struc) a() int { return s.x }
func (s *struc) b() int { return s.y }

func use(s *struc, c bool) {
	_ = s.a()
	_ = s.b()
	f := s.a
	if c {
		f = s.b
	}
	print(f())
}

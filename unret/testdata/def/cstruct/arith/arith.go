package csarith

type struc struct {
	a, b, c int
}

func (s *struc) A() int { return s.a }
func (s *struc) B() int { return s.b }
func (s *struc) C() int { return s.c }

func use() {
	s := struc{3, 4, 5}
	if s.A()*s.A()+s.B()*s.B() == s.C()*s.C() {
		println("right")
	}
}

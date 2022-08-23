package csclosure

type struc struct{ a, b, c int }

func (s *struc) val1() int { return s.a }
func (s *struc) val2() int { return s.b }
func (s *struc) val3() int { return s.c }

func use() {
	s := struc{3, 4, 5}

	x := s.val1()
	y := func() int {
		return x + s.val2()
	}
	z := func() int {
		return y() + s.val3()
	}()

	println(z)
}

package sclosure

type struc struct{ a, b, c int }

func (s *struc) vals() (a, b, c int) { // want "2 \\(c int\\) is never used"
	return s.a, s.b, s.c
}

func use() {
	s := struc{3, 4, 5}

	f := func() int { // want "0 \\(int\\) is never used"
		x, y, _ := s.vals()
		return x + y
	}
	_ = f
}

package csassign

var out struct {
	FuncInt func() int
}

type struc struct{ a int }

func (s *struc) val1() int { return s.a }
func (s *struc) val2() int { return s.a }

func use(s *struc) {
	_ = s.val1()
	out.FuncInt = s.val1

	_ = s.val2()
	made := struct {
		FuncInt func() int
	}{
		FuncInt: s.val2,
	}
	_ = made
}

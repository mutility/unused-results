package csretval

type struc struct{}

func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use() (int, string, error) {
	var s struc
	return s.rets3()
}

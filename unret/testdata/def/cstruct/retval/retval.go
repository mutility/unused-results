package csretvaltup

type struc struct{}

func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use() (string, int, error) {
	var s struc
	a, b, c := s.rets3()
	return b, a, c
}

package recvclose

type struc struct{}

// TODO:evaluate; should probably be 0 or 3 diagnostics on rets3, not 2.

func (*struc) rets0()                      {}
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // want "never" "never"

func use(s *struc) {
	use0 := func(fn func()) {}
	use3 := func(fn func() (int, string, error)) {}

	use0(s.rets0)
	use0(s.rets0)

	// multiple instances of s.rets3 generate duplicate closures

	use3(s.rets3)
	use3(s.rets3)
}

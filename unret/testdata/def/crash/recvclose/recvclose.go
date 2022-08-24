package recvclose

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // want "0 .int. is never used" "1 .string. is never used" "2 .error. is never used"

func use(s *struc) {
	use0 := func(fn func()) {}
	use3 := func(fn func() (int, string, error)) {}

	use0(s.rets0)
	use0(s.rets0)

	// multiple instances of s.rets3 generate duplicate closures

	use3(s.rets3)
	use3(s.rets3)
}

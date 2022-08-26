package stypecall

type f func() error

type struc struct{}

func (*struc) f() error { return nil } // want "0 \\(error\\) is never used"

func use(s *struc, ff func(f)) {
	f(s.f)()
}

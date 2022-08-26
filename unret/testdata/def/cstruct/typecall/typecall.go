package cstypecall

type f func() error

type struc struct{}

func (*struc) f() error { return nil }

func use(s *struc, ff func(f)) {
	ff(s.f)
}

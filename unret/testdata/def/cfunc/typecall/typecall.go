package cftypecall

type f func() error

func use(ff func(f)) {
	fn := func() error {
		return nil
	}
	_ = fn()
	ff(fn)
}

package ftypecall

type f func() error

func use(ff func(f)) {
	fn := func() error { // want "0 \\(error\\) is never used"
		return nil
	}

	x := f(fn)
	x()
}

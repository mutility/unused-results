package cfcond

func a() int { return 3 }
func b() int { return 4 }

func use(c bool) {
	_ = a()
	_ = b()
	f := a
	if c {
		f = b
	}
	print(f())
}

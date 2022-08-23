package ideferred

type itf interface {
	vals() (a, b, c int) // want "1 \\(b int\\) is never used"
}

func use(i itf) {
	x, _, _ := i.vals()
	defer func(int) {
		x, _, z := i.vals()
		println(x + z)
	}(x)
}

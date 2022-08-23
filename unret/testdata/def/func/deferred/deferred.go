package fdeferred

func vals() (a, b, c int) { return 3, 4, 5 } // want "1 \\(b int\\) is never used"

func use() {
	x, _, _ := vals()
	defer func(int) {
		x, _, z := vals()
		println(x + z)
	}(x)
}

package cfdeferred

func vals() (a, b, c int) { return 3, 4, 5 }

func use() {
	x, _, _ := vals()
	defer func(int) {
		_, y, z := vals()
		println(y + z)
	}(x)
}

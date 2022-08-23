package cideferred

type itf interface {
	vals() (a, b, c int)
}

func use(i itf) {
	x, _, _ := i.vals()
	defer func(int) {
		_, y, z := i.vals()
		println(y + z)
	}(x)
}

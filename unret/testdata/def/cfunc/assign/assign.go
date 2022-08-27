package cfassign

var out struct {
	FuncInt func() int
}

func val1() int { return 3 }
func val2() int { return 3 }

func use() {
	_ = val1()
	out.FuncInt = val1

	_ = val2()
	made := struct {
		FuncInt func() int
	}{
		FuncInt: val2,
	}
	_ = made
}

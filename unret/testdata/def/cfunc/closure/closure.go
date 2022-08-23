package cfclosure

func val1() int { return 3 }
func val2() int { return 4 }
func val3() int { return 5 }

func use() {
	x := val1()
	y := func() int {
		return x + val2()
	}
	z := func() int {
		return y() + val3()
	}()

	println(z)
}

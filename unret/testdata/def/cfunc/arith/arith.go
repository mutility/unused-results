package cfarith

func val1() int { return 3 }
func val2() int { return 4 }
func val3() int { return 5 }

func use() {
	if val1()*val1()+val2()*val2() == val3()*val3() {
		println("right")
	}
}

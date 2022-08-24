package cfexported

import "log"

// Since these functions are exported, they are not warned by default.
func Res0()                      {}
func Res1() int                  { return 1 }
func Res2() (x, y int)           { return 2, 3 }
func Res3() (int, string, error) { return 4, "5", nil }

func use() {
	_ = Res1()
	a, _ := Res2()
	_, b, _ := Res3()
	_, _, c := Res3()

	log.Print(a, b, c)
}

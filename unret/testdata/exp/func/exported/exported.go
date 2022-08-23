package xfexported

import "log"

// Since these functions are exported, they are not warned by default.

func Res0()
func Res1() int                  // want "result 0 \\(int\\) is never used"
func Res2() (x, y int)           // want "result 1 \\(y int\\) is never used"
func Res3() (int, string, error) // want "result 0 \\(int\\) is never used"

func use() {
	Res0()
	Res1()
	a, _ := Res2()
	_, b, _ := Res3()
	_, _, c := Res3()

	log.Print(a, b, c)
}

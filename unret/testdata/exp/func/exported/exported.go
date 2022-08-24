package xfexported

import "log"

// Since these functions are exported, they are not warned by default.

func Res0()                      {}
func Res1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func Res2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func Res3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use() {
	Res0()
	Res1()
	a, _ := Res2()
	_, b, _ := Res3()
	_, _, c := Res3()

	log.Print(a, b, c)
}

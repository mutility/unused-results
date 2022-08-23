package fprinted

import "log"

func rets0()                      {}
func rets1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func rets2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func rets3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use() {
	rets0()
	rets1()
	a, _ := rets2()
	_, b, _ := rets3()
	_, _, c := rets3()

	log.Print(a, b, c)
}

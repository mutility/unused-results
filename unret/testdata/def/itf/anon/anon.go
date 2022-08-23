package ianon

import "log"

// TODO: Since the interface is referenced, we should report on it.

func use(i interface {
	rets0()
	rets1() int                  // TODO:want "result 0 \\(int\\) is never used"
	rets2() (x, y int)           // TODO:want "result 1 \\(y int\\) is never used"
	rets3() (int, string, error) // TODO:want "result 0 \\(int\\) is never used"
},
) {
	i.rets1()
	a, _ := i.rets2()
	_, b, _ := i.rets3()
	_, _, c := i.rets3()

	log.Print(a, b, c)
}

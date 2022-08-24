package sanonitf

import "log"

// Since struc{} is the static implementation, it can be reported upon directly.

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }           // want "result 0 \\(int\\) is never used"
func (*struc) rets2() (x, y int)           { return 2, 3 }        // want "result 1 \\(y int\\) is never used"
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used"

func use() {
	i := interface {
		rets0()
		rets1() int                  // want "result 0 \\(int\\) is never used"
		rets2() (x, y int)           // want "result 1 \\(y int\\) is never used"
		rets3() (int, string, error) // want "result 0 \\(int\\) is never used"
	}(&struc{})
	i.rets0()
	_ = i.rets1()
	a, _ := i.rets2()
	_, b, _ := i.rets3()
	_, _, c := i.rets3()

	log.Print(a, b, c)
}

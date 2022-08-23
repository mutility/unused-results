package csmakeitf

import "log"

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

type itf interface {
	rets0()
	rets1() int
	rets2() (x, y int)
	rets3() (int, string, error)
}

func use() {
	i := itf(&struc{})
	i.rets0()
	log.Print(i.rets1())
	log.Print(i.rets2())
	log.Print(i.rets3())
}

package ciprinted

import "log"

type itf interface {
	rets0()
	rets1() int
	rets2() (x, y int)
	rets3() (int, string, error)
}

func use(i itf) {
	i.rets0()
	log.Print(i.rets1())
	log.Print(i.rets2())
	log.Print(i.rets3())
}

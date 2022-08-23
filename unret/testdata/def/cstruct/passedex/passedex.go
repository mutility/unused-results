package cspassedex

import "log"

type struc struct{}

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }
func (*struc) rets2() (x, y int)           { return 2, 3 }
func (*struc) rets3() (int, string, error) { return 4, "5", nil }

func use(s *struc) {
	log.Print(s.rets0)
	// TODO: log.Print(s.rets1)
	// TODO: log.Print(s.rets2)
	// TODO: log.Print(s.rets3)
}

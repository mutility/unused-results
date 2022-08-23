package cfpassedex

import "log"

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }
func rets3() (int, string, error) { return 4, "5", nil }

func use() {
	log.Print(rets0)
	log.Print(rets1)
	log.Print(rets2)
	log.Print(rets3)
}

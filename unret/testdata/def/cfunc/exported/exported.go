package cfexported

import "log"

// Since these functions are exported, they are not warned by default.

func NotAllUsed1() int
func NotAllUsed2() (x, y int)
func NotAllUsed3() (int, string, error)

func use() {
	_ = NotAllUsed1()
	a, _ := NotAllUsed2()
	_, b, _ := NotAllUsed3()
	_, _, c := NotAllUsed3()

	log.Print(a, b, c)
}

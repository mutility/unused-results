package funcalled

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }        // want "result 0 \\(x int\\) is never used" "result 1 \\(y int\\) is never used"
func rets3() (int, string, error) { return 4, "5", nil } // want "result 0 \\(int\\) is never used" "result 1 \\(string\\) is never used" "result 2 \\(error\\) is never used"

func use() {
	println(rets1())
}

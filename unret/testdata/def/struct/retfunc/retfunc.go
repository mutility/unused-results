package sretfunc

type struc struct{}

/* TODO
// since neither struc nor use are exported, this does not count as use

func (*struc) rets0()                      {}
func (*struc) rets1() int                  { return 1 }           // TODO:want "result 0 \\(int\\) is never used"
func (*struc) rets2() (x, y int)           { return 2, 3 }        // TODO:want "result 0 \\(x int\\) is never used" "result 1 \\(y int\\) is never used"
func (*struc) rets3() (int, string, error) { return 4, "5", nil } // TODO:want "result 0 \\(int\\) is never used" "1 \\(string\\) is never used" "2 \\(error\\) is never used"

func use() (func(), func() int, func() (int, int), func() (int, string, error)) {
	var s struc
	return s.rets0, s.rets1, s.rets2, s.rets3
}
*/

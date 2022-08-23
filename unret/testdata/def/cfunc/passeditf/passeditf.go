package cfpasseditf

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }
func rets3() (int, string, error) { return 4, "5", nil }

func use(itf interface{ use(...interface{}) }) {
	itf.use(rets0)
	itf.use(rets1)
	itf.use(rets2)
	itf.use(rets3)
}

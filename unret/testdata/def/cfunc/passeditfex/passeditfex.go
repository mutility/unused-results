package cfpasseditfex

func rets0()                      {}
func rets1() int                  { return 1 }
func rets2() (x, y int)           { return 2, 3 }
func rets3() (int, string, error) { return 4, "5", nil }

func use(itf interface{ Use(...interface{}) }) {
	itf.Use(rets0)
	itf.Use(rets1)
	itf.Use(rets2)
	itf.Use(rets3)
}

package f

type foo struct{}

func (*foo) many() (a, b, c int) { return 3, 4, 5 } // want "c int\\) is never used"

func useNamed(f *foo) (sum int) {
	defer func() {
		a, b, _ := f.many()
		sum = a + b
	}()
	return 0
}

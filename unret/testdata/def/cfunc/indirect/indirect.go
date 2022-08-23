package cfindirect

var glob int

func sink() *int  { return &glob }
func source() int { return 4 }

func use() {
	*sink() = source()
}

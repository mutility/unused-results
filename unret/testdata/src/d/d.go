package d

import "log"

func foo() {
	bar := func() (int, error) { // TODO:want "error is never used"
		return 1, nil
	}
	n, _ := bar()
	log.Println(n)
}

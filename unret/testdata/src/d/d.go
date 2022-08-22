package d

import "log"

func foo() {
	bar := func() (int, error) { // want "error\\) is never used"
		return 1, nil
	}
	n, _ := bar()
	log.Println(n)

	baz := func() (int, error) { // want "error\\) is never used"
		return n + 1, nil
	}
	n, _ = baz()
	log.Println(n)
}

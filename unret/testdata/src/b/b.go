package b

import "log"

type foo struct{}

func (*foo) used() error   { return nil }
func (*foo) unused() error { return nil } // want "error is never used"

type bar struct{}

func (*bar) used() error   { return nil }
func (*bar) unused() error { return nil } // want "error is never used"

func use() {
	f := foo{}
	log.Println(f.used())
	_ = f.unused()

	type itf interface {
		used() error
		unused() error
	}

	i := itf(&bar{})
	log.Println(i.used())
	_ = i.unused()
}

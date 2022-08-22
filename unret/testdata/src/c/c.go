package c

import "log"

type foo struct{}

func (*foo) used() error   { return nil }
func (*foo) unused() error { return nil }

type bar struct{}

func (*bar) used() error   { return nil }
func (*bar) unused() error { return nil }

type itf interface {
	used() error
	unused() error // want "error is never used"
}

func useNamed(i itf) {
	log.Println(i.used())
	_ = i.unused()
}

func useAnon(i interface {
	used() error
	unused() error // TODO:want "error is never used"
},
) {
	log.Println(i.used())
	_ = i.unused()
}

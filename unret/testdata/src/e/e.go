package e

type itf1 interface {
	read() (int, error) // want "int\\) is never used"
}

func deferred(i itf1) (err error) {
	defer func() {
		_, err = i.read()
	}()
	return
}

type itf2 interface {
	read() (int, error) // want "int\\) is never used"
}

func deferredPassItf(i itf2) (err error) {
	defer func(i itf2) {
		_, err = i.read()
	}(i)
	return
}

type itf3 interface {
	read() (int, error) // want "int\\) is never used"
}

func deferredPassErr(i itf3) (err error) {
	defer func(perr *error) {
		_, *perr = i.read()
	}(&err)
	return
}

type itf4 interface {
	read() (int, error) // want "int\\) is never used"
}

func deferredPassBoth(i itf4) (err error) {
	defer func(i itf4, p *error) {
		_, *p = i.read()
	}(i, &err)
	return
}

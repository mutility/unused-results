package g

import "log"

type closed struct {
	data int
}

func (c *closed) foo() int                  { return c.data }
func (c *closed) returnedIn() (int, error)  { return c.data, nil } // TODO:want "int\\) is never used"
func (c *closed) returnedOut() (int, error) { return c.data, nil }
func (c *closed) used() (int, error)        { return c.data, nil }
func (c *closed) passedIn() (int, error)    { return c.data, nil } // TODO:want "int\\) is never used"
func (c *closed) passedOut() (int, error)   { return c.data, nil }

func bar() {
	get := func(fn func() int) int { return fn() } // want "int\\) is never used"
	cl := closed{}
	get(cl.foo) // generates duplicate closures.
	get(cl.foo) // Both synthetic.
}

// retsfunc is not exported, but returns a closure.
// don't consider it used by default.
func retsfunc() func() (int, error) {
	cl := closed{}
	_, err := cl.returnedIn()
	log.Print(err)
	return cl.returnedIn
}

// RetsFunc is exported, and returns a closure.
// Consider it all used by default.
func RetsFunc() func() (int, error) {
	cl := closed{}
	_, err := cl.returnedOut()
	log.Print(err)
	return cl.returnedOut
}

func PassesToItfIn(itf interface{ use(func() (int, error)) }) {
	cl := closed{}
	_, err := cl.passedIn()
	log.Print(err)
	itf.use(cl.passedIn)
}

func PassesToItfOut(itf interface{ Use(func() (int, error)) }) {
	cl := closed{}
	_, err := cl.passedOut()
	log.Print(err)
	itf.Use(cl.passedOut)
}

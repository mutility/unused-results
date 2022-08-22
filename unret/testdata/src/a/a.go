package a

import "log"

func norets()                     {}
func usedret() error              { return nil }
func usedrets() (string, error)   { return "hi", nil }
func unusedret() error            { return nil }       // want "error is never used"
func unusedret1() (string, error) { return "hi", nil } // want "string is never used"
func unusedret2() (string, error) { return "hi", nil } // want "error is never used"
func unusedrets() (string, error) { return "hi", nil } // want "string is never used" "error is never used"
func unusedfunc() (string, error) { return "hi", nil }
func hello() string               { return "hi" }

func Norets()                     {}
func Usedret() error              { return nil }
func Usedrets() (string, error)   { return "hi", nil }
func Unusedret() error            { return nil }
func Unusedret1() (string, error) { return "hi", nil }
func Unusedret2() (string, error) { return "hi", nil }
func Unusedrets() (string, error) { return "hi", nil }
func Unusedfunc() (string, error) { return "hi", nil }

func uses() {
	norets()
	if err := usedret(); err != nil {
		log.Println("Error!", err)
	}
	if s, err := usedrets(); err != nil {
		log.Println("String", s)
	}
	_ = unusedret()
	if _, err := unusedret1(); err != nil {
		log.Println("Error!", err)
	}
	if s, _ := unusedret2(); s != "" {
		log.Println("String", s)
	}
	_, _ = unusedrets()
	log.Println(hello())
}

func usesExported() {
	Norets()
	if err := Usedret(); err != nil {
		log.Println("Error!", err)
	}
	if s, err := Usedrets(); err != nil {
		log.Println("String", s)
	}
	_ = Unusedret()
	if _, err := Unusedret1(); err != nil {
		log.Println("Error!", err)
	}
	if s, _ := Unusedret2(); s != "" {
		log.Println("String", s)
	}
	_, _ = Unusedrets()
}

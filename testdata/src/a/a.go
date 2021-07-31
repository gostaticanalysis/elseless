package a

func f() {
	if 1 == 1 {
		println("1")
		return
	} else if 2 == 2 { // want "unnecessary else"
		println("2")
	} else { // OK
		println("3")
	}
}

package a

func f() {
	if true {
		println("1")
		return
	} else if true { // want "unnecessary else"
		println("2")
	} else { // OK
		println("3")
	}

	for {
		if true {
			continue
		} else { // want "unnecessary else"
			println()
		}

		if true {
			break
		} else { // want "unnecessary else"
			println()
		}

		switch {
		case true:
			if true {
				break
			} else { // want "unnecessary else"
				println()
			}
		}
	}

	if true {
		panic(nil)
	} else { // want "unnecessary else"
		println()
	}
}

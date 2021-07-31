# elseless

[![pkg.go.dev][gopkg-badge]][gopkg]

`elseless` finds unnecessary `else`.

```go
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
```

`fixelseless` command check and fix (remove) unnecessary `else`.

```sh
$ go install github.com/gostaticanalysis/elseless/cmd/fixelseless@latest
$ fixelseless ./...
```
<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/elseless
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/elseless?status.svg


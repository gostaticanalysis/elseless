package main

import (
	"github.com/gostaticanalysis/elseless"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(elseless.Analyzer) }

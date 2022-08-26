package main

import (
	"github.com/mutility/analyzers/unret"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(unret.Analyzer().Analyzer)
}

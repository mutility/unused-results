// Package main hosts the unused-results analyzer unret.Analyzer()
package main

import (
	"github.com/mutility/unused-results/v1/unret"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(unret.Analyzer().Analyzer)
}

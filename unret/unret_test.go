package unret_test

import (
	"testing"

	"github.com/mutility/analyzers/unret"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, unret.Analyzer, "a", "b", "c", "d", "e", "f")
}

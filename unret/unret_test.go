package unret_test

import (
	"testing"

	"github.com/mutility/analyzers/unret"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	t.Run("defaults", func(t *testing.T) {
		analysistest.Run(t, testdata, unret.Analyzer, "./def/...")
	})
	t.Run("exported", func(t *testing.T) {
		unret.Analyzer.Flags.Lookup("exported").Value.Set("true")
		analysistest.Run(t, testdata, unret.Analyzer, "./exp/...")
		unret.Analyzer.Flags.Lookup("exported").Value.Set("false")
	})
	t.Run("uncalled", func(t *testing.T) {
		unret.Analyzer.Flags.Lookup("uncalled").Value.Set("true")
		analysistest.Run(t, testdata, unret.Analyzer, "./unc/...")
		unret.Analyzer.Flags.Lookup("uncalled").Value.Set("false")
	})
	t.Run("passed", func(t *testing.T) {
		unret.Analyzer.Flags.Lookup("passed").Value.Set("true")
		analysistest.Run(t, testdata, unret.Analyzer, "./pas/...")
		unret.Analyzer.Flags.Lookup("passed").Value.Set("false")
	})
}

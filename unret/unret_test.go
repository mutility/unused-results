package unret_test

import (
	"testing"

	"github.com/mutility/analyzers/unret"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	t.Run("defaults", func(t *testing.T) {
		t.Parallel()
		analysistest.Run(t, testdata, unret.Analyzer().Analyzer, "./def/...")
	})
	t.Run("exported", func(t *testing.T) {
		t.Parallel()
		u := unret.Analyzer()
		u.ReportExported = true
		analysistest.Run(t, testdata, u.Analyzer, "./exp/...")
	})
	t.Run("uncalled", func(t *testing.T) {
		t.Parallel()
		u := unret.Analyzer()
		u.ReportUncalled = true
		analysistest.Run(t, testdata, u.Analyzer, "./unc/...")
	})
	t.Run("passed", func(t *testing.T) {
		t.Parallel()
		u := unret.Analyzer()
		u.ReportPassed = true
		analysistest.Run(t, testdata, u.Analyzer, "./pas/...")
	})
	t.Run("returned", func(t *testing.T) {
		t.Parallel()
		u := unret.Analyzer()
		u.ReportReturned = true
		analysistest.Run(t, testdata, u.Analyzer, "./ret/...")
	})
}

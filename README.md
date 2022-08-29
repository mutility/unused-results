# Unused Results

`unused-results` is a
[go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)-based tool
that identifies ignored results from the functions, interfaces, and closures in
your package. It is meant as a tool to help you find flaws in your design or
implementation. Just remember that not all unused results are flaws.

## Example messages

Given the following source code `example.go` simplified from
[anonitf.go](unret/testdata/struct/anonitf/anonitf.go):

```go
     1	package example
     2	
     3	type struc struct{}
     4	
     5	func (*struc) rets2() (x, y int) { return 2, 3 }
     6	
     7	func use(s *struc) {
     8	    i := interface{ rets2() (int, int) }(s)
     9	    a, _ := i.rets2()
    10	    println(a)
    11	}
```

unused-results will report the following. Note that in cases like these it
identifies both the interface and implementation whose results are unused.
 
```console
$t unused-results ./...
.../example.go:5:15: rets2 result 1 (y int) is never used
.../example.go:8:21: rets2 result 1 (int) is never used
exit status 3
```

## Usage

Run from source with `go run github.com/mutility/unused-results@latest` or
install with `go install github.com/mutility/unused-results@latest` and run
unused-results from GOPATH/bin.

You can configure behvior at the command line by passing the flags below, or in
library use by setting fields on `unret.Analyzer()`. All of these default to
false.

Flag | Field | Meaning
-|-|-
`-exported` | ReportExported | Allow exported functions to be reported
`-uncalled` | ReportUncalled | Allow uncalled functions to be reported
`-passed` | ReportPassed | Allow functions passed to others to be reported
`-returned` | ReportReturned | Allow functions returned from others to be reported
`-assigned` | ReportAssigned | Allow functions assigned to storage to be reported

Notes:

- Reports of a function can be omitted for multiple reasons; all relevant flags
  or fields must be set to reveal such an unused result
- References to a function or closure object are often considered as a call to
  that function

## False positives and negatives

Due to the package-by-package nature of how analysis-based linters work, it is
not feasible to consider external uses of exported or returned functions. It
may be possible to consider use patterns of external functions that are called
with one of the local functions, but this is not currently implemented.

To avoid too many false positives there are unfortunately also a lot of
potentially false negatives. These are the defaults:

- Exported functions are assumed to be fully used
- Uncalled functions are not reported
- Functions passed as closures to other functions are not reported
- Functions returned as closures from other functions are not reported
- Functions assigned to storage such as fields, globals, etc. are not reported

On the other hand, if you have a simple unused result that you don't want to
see, you cannot silence the report by assigning to `_`, nor by assigning to a
local that is later only assigned to `_`. You must do something more involved
to fool `unused-results`, such as storing in an anonymous field of a discarded
struct, or passing to another function.

```go
    ..., otherwiseUnused, ... = yourFunc()
    var _ = struct{any}{otherwiseUnused}
    // or 
    func(...any){} (otherwiseUnused)
```

Note that if tracing is improved, the simple cases here may become insufficient
to fool `unused-results`.

## Differences from govet unusedResults

[`govet`](https://pkg.go.dev/cmd/vet) includes a check called unusedresult.
This check looks for callers that are ignoring important results from functions
they call.

By contrast, unused-results looks at each each function, interface, or
closure and attempts to identify the results that might be unimportant because
no visible caller uses them.

## Bug reports and feature contributions

`unused-results` is developed in spare time, so while bug reports and feature
contributions are welcomed, it may take a while for them to be reviewed. If
possible, try to find a minimal reproduction before reporting a bug. Bugs that
are difficult or impossible to reproduce will likely be closed.

All bug fixes will include tests to help ensure no regression; correspondingly
all contributions should include such tests.

## Mutility Analyzers

`unused-results` is part of [mutility-analyzers](https://github.com/mutility/analyzers).

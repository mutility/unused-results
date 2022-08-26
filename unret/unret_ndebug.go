//go:build !debug

package unret

import (
	"go/token"
	"log"

	"golang.org/x/tools/go/ssa"
)

func debugDumpFuncSSA(fn *ssa.Function, debug *log.Logger)                                  {}
func debugDumpUses(debug *log.Logger, fset *token.FileSet, used map[callee]usage)           {}
func debugTrackApply(debug *log.Logger, inst ssa.Instruction, fun callee, set func(*usage)) {}

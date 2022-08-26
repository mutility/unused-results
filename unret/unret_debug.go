//go:build debug

package unret

import (
	"go/token"
	"go/types"
	"io"
	"log"

	"golang.org/x/tools/go/ssa"
)

func debugDumpFuncSSA(fn *ssa.Function, debug *log.Logger) {
	if len(fn.Blocks[0].Instrs) > 1 && debug.Writer() != io.Discard {
		flags := debug.Flags()
		defer debug.SetFlags(flags)
		debug.SetFlags(0)
		for _, blk := range fn.Blocks {
			debug.Println(fn, blk)
			for _, inst := range blk.Instrs {
				name, typ := "", ""
				if val, ok := inst.(ssa.Value); ok {
					name = val.Name() + " ="
					typ = val.Type().String()
				}
				debug.Printf("%7.7s %-30.30s %20.20s [%[2]T]", name, inst, typ)
			}
		}
	}
}

func debugDumpUses(debug *log.Logger, fset *token.FileSet, used map[callee]usage) {
	debug.Println("USES:", used)
	for fun, res := range used {
		pos := fset.Position(fun.Pos())
		nres := fun.Type().(*types.Signature).Results().Len()
		if nres > 0 {
			debug.Printf("  %s:%d: %s() %x/%x passed:%t returned:%t", pos.Filename, pos.Line, fun.Name(), res.results, uint(1<<nres)-1, res.passed, res.returned)
		}
	}
}

func debugTrackApply(debug *log.Logger, inst ssa.Instruction, fun callee, set func(*usage)) {
	var val string
	if v, ok := inst.(ssa.Value); ok {
		val = v.Name()
	}
	direct := usage{}
	set(&direct)
	debug.Println("USED", fun.Name(), direct, "by", inst, val)
}

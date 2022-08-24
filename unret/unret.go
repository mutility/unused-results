package unret

import (
	"go/token"
	"go/types"
	"io"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = `unret reports returns from unexported functions that are never used.`

var reportExported, reportUncalled, debugAnalyzer bool

func init() {
	Analyzer.Flags.BoolVar(&reportExported, "exported", false, "report unused results from exported functions")
	Analyzer.Flags.BoolVar(&reportUncalled, "uncalled", false, "report unused results from uncalled functions")
	Analyzer.Flags.BoolVar(&debugAnalyzer, "verbose", true, "issue debug logging")
}

var Analyzer = &analysis.Analyzer{
	Name:     "unret",
	Doc:      doc,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	type poser interface {
		Name() string
		Pos() token.Pos
		Type() types.Type
	}

	debug := log.New(log.Default().Writer(), pass.Pkg.Name()+": ", log.Lshortfile)
	if !debugAnalyzer {
		debug.SetOutput(io.Discard)
	}
	type usage struct {
		results uint // track explicitly used results as a bit field
	}
	prog := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	used := make(map[poser]usage, len(prog.SrcFuncs))

	recvType := func(t types.Type) *types.TypeName {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		}
		if nam, ok := t.(*types.Named); ok && nam.Obj() != nil {
			return nam.Obj()
		}
		return nil
	}

	// typeFuncs stores typename->funcname->poser so that methods used
	// via interface can be tracked. See testdata/src/b.
	funcs := make([]poser, 0, len(prog.SrcFuncs))
	typeFuncs := make(map[*types.TypeName]map[string]poser)
	for _, fun := range prog.SrcFuncs {
		funcs = append(funcs, fun)
		if recv := fun.Signature.Recv(); recv != nil {
			if typ := recvType(recv.Type()); typ != nil {
				m, ok := typeFuncs[typ]
				if !ok {
					m = make(map[string]poser)
					typeFuncs[typ] = m
				}
				m[fun.Name()] = fun
			}
		}
	}

	for _, name := range pass.Pkg.Scope().Names() {
		obj := pass.Pkg.Scope().Lookup(name)
		typ := obj.Type()
		if nam, ok := typ.(*types.Named); ok {
			typ = nam.Underlying()
		}
		if itf, ok := typ.(*types.Interface); ok {
			m, ok := typeFuncs[obj.(*types.TypeName)]
			if !ok {
				m = make(map[string]poser, itf.NumMethods())
				typeFuncs[obj.(*types.TypeName)] = m
			}
			for i, n := 0, itf.NumMethods(); i < n; i++ {
				fn := itf.Method(i)
				m[fn.Name()] = fn
				funcs = append(funcs, fn)
			}
		}
	}

	// for k, v := range typeFuncs {
	// 	debug.Println("tf:", k, v)
	// }
	anonTypes := make(map[string]*types.TypeName)  // track anon interfaces
	returned := make(map[*types.TypeName]struct{}) // track types that get returned
	closures := make(map[*ssa.Function]*ssa.MakeClosure)
	for _, fun := range prog.SrcFuncs {
		if token.IsExported(fun.Name()) {
			res := fun.Signature.Results()
			for i, nr := 0, res.Len(); i < nr; i++ {
				if t := recvType(res.At(i).Type()); t != nil {
					returned[t] = struct{}{}
				}
			}
		}
		for _, block := range fun.Blocks {
			for _, instr := range block.Instrs {
				if cl, ok := instr.(*ssa.MakeClosure); ok && len(cl.Bindings) > 0 {
					fn := cl.Fn.(*ssa.Function)
					if fn.Synthetic != "" {
						continue // ignore these; see testdata/src/g
					}
					if prev, ok := closures[fn]; ok {
						debug.Println("Repeat closure", prev, cl, pass.Fset.Position(prev.Pos()), pass.Fset.Position(cl.Pos()))
					}
					closures[fn] = cl
				}
			}
		}
	}

	trackAnonInterface := func(itf *types.Interface) *types.TypeName {
		anon, ok := anonTypes[itf.String()]
		if !ok {
			anon = types.NewTypeName(token.NoPos, prog.Pkg.Pkg, itf.String(), itf)
			anonTypes[itf.String()] = anon
			m := make(map[string]poser, itf.NumExplicitMethods())
			typeFuncs[anon] = m
			for i, ni := 0, itf.NumExplicitMethods(); i < ni; i++ {
				meth := itf.Method(i)
				m[meth.Name()] = meth
				funcs = append(funcs, meth)
			}
		}
		return anon
	}

	// funResult returns the poser and result index that a ssa.Value
	// uses. Most of these indicate use of that result, but *ssa.Extract isn't
	// itself a use.
	var funResult func(op ssa.Value, extract func(*types.TypeName) poser) (res []poser, idx int)
	funResult = func(op ssa.Value, extract func(*types.TypeName) poser) (res []poser, idx int) {
		switch op := op.(type) {
		case *ssa.Extract:
			fun, _ := funResult(op.Tuple, extract)
			return fun, op.Index
		case *ssa.Call:
			com := op.Common()
			if sfunc := com.StaticCallee(); sfunc != nil {
				return append(res, sfunc), 0
			}
			if com.Method != nil {
				extract = func(tn *types.TypeName) poser {
					return typeFuncs[tn][com.Method.Name()]
				}
			}
			f, i := funResult(com.Value, extract)
			return f, i
		case *ssa.MakeInterface:
			f, i := funResult(op.X, extract)
			switch t := op.Type().(type) {
			case *types.Interface:
				return append(f, extract(trackAnonInterface(t))), i
			case *types.Named:
				if f2 := extract(t.Obj()); f2 != nil {
					return append(f, f2), i
				}
			}
			return f, i
		case *ssa.Alloc:
			if typ := recvType(op.Type().Underlying().(*types.Pointer).Elem()); typ != nil {
				return append(res, extract(typ)), 0
			}
			return nil, 0
		case *ssa.Parameter:
			if typ := recvType(op.Type()); typ != nil {
				ex := extract(typ)
				return append(res, ex), 0
			}
			if itf, ok := op.Type().(*types.Interface); ok {
				return append(res, extract(trackAnonInterface(itf))), 0
			}
			return nil, 0
		case *ssa.UnOp:
			return funResult(op.X, extract)
		case *ssa.FreeVar:
			for i, fv := range op.Parent().FreeVars {
				if fv != op {
					continue
				}
				return funResult(closures[op.Parent()].Bindings[i], extract)
			}
			return nil, 0

		case *ssa.MakeClosure:
			if sfn := op.Fn.(*ssa.Function); sfn.Synthetic != "" {
				if tfn, ok := sfn.Object().(*types.Func); ok {
					sig := tfn.Type().(*types.Signature)
					return append(res, typeFuncs[recvType(sig.Recv().Type())][tfn.Name()]), 0
				}
			}
			return append(res, op.Fn), 0

		case *ssa.ChangeInterface:
			// debug.Println("ChangeInterface:", op.Type(), reflect.TypeOf(op.Type()), op.X, reflect.TypeOf(op.X))
			return funResult(op.X, extract)
		default:
			// debug.Printf("whattabout %T: %v", op, op)
			return nil, 0
		}
	}

	nilExtract := func(*types.TypeName) poser { return nil }
	for _, fn := range prog.SrcFuncs {
		// dumpfunc(fn, debug)
		for _, blk := range fn.Blocks {
			for _, inst := range blk.Instrs {
				if val, ok := inst.(ssa.Value); ok {
					switch val := val.(type) {
					// consider extract as function uses, but not as result uses
					case *ssa.Call:
						funs, _ := funResult(val, nilExtract)
						for _, fun := range funs {
							if r, ok := used[fun]; fun != nil && !ok {
								used[fun] = r
							}
						}
					case *ssa.Extract:
						funs, _ := funResult(val, nilExtract)
						for _, fun := range funs {
							if r, ok := used[fun]; fun != nil && !ok {
								used[fun] = r
							}
						}
						continue // extracts are only a function use, so skip its operands
					}
				}

				// apply updates records, both direct and any indirect ones.
				// inst is used only for (disabled) logging.
				apply := func(inst ssa.Instruction, fun poser, set func(*usage)) {
					r := used[fun]
					set(&r)
					used[fun] = r
					if sfunc, ok := fun.(*ssa.Function); ok {
						if fun, ok := sfunc.Object().(*types.Func); ok {
							r := used[fun]
							set(&r)
							used[fun] = r
							// var val string
							// if v, ok := inst.(ssa.Value); ok {
							// 	val = v.Name()
							// }
							// debug.Println("USED", fun.Name(), up(returns{}), "by", inst, val)
						}
					}
				}

				// consider other instructions referring to call/extract as function and result uses
				for _, op := range inst.Operands(nil) {
					if op == nil || *op == nil {
						continue
					}
					funs, res := funResult(*op, nilExtract)
					for _, fun := range funs {
						apply(inst, fun, func(r *usage) { r.results |= 1 << res })
					}
				}
			}
		}
	}

	// debug.Println("USES:", used)
	// for fun, res := range used {
	// 	pos := pass.Fset.Position(fun.Pos())
	// 	nres := fun.Type().(*types.Signature).Results().Len()
	// 	if nres > 0 {
	// 		debug.Printf("  %s:%d: %s() %x/%x", pos.Filename, pos.Line, fun.Name(), res.used, uint(1<<nres)-1)
	// 	}
	// }

	for _, fn := range funcs {
		if !reportExported && token.IsExported(fn.Name()) {
			switch fn := fn.(type) {
			case *ssa.Function:
				if fn.Signature.Recv() == nil {
					continue
				}

				t := recvType(fn.Signature.Recv().Type())
				if t == nil || t.Exported() {
					continue
				}
				if _, ok := returned[t]; ok {
					continue
				}
			default:
				continue
			}
		}
		r, ok := used[fn]
		if !reportUncalled && !ok {
			continue
		}
		results := fn.Type().(*types.Signature).Results()
		for i, n := 0, results.Len(); i < n; i++ {
			if (1<<uint(i))&r.results == 0 {
				res := results.At(i)
				rnt := strings.TrimSpace(res.Name() + " " + res.Type().String())
				pass.Reportf(fn.Pos(), "%s result %d (%s) is never used", fn.Name(), i, rnt)
				// debug.Println(fn.Name(), "result", i, "unused")
			}
		}
	}

	return nil, nil
}

func dumpfunc(fn *ssa.Function, debug *log.Logger) {
	if len(fn.Blocks[0].Instrs) > 4 && debug.Writer() != io.Discard {
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

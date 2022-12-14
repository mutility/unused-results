// Package unret implements the unused-results analyzer. It reports per callee
// (a function, interface, or closure) the results that were ignored. By
// default it does not report results from exported or uncalled callees, nor
// from ones that appear likely to be exported in some other way. These
// defaults can be overridden.
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

// unretAnalyzer offers configuration options for the unused-results analyzer.
type unretAnalyzer struct {
	*analysis.Analyzer
	ReportExported bool // report unused results from exported functions
	ReportUncalled bool // report unused results from uncalled functions
	ReportPassed   bool // report unused results from functions passed to other functions
	ReportReturned bool // report unused results from functions returned by other
	ReportAssigned bool // report unused results from functions assigned to storage
	Debug          bool
}

// Analyzer returns a new unretAnalyzer that can be configured before using
// with, e.g., an analysis checker or test harness.
func Analyzer() *unretAnalyzer {
	u := &unretAnalyzer{
		Analyzer: &analysis.Analyzer{
			Name:     "unret",
			Doc:      doc,
			Requires: []*analysis.Analyzer{buildssa.Analyzer},
		},
	}
	u.Flags.BoolVar(&u.ReportExported, "exported", false, "report unused results from exported functions")
	u.Flags.BoolVar(&u.ReportUncalled, "uncalled", false, "report unused results from uncalled functions")
	u.Flags.BoolVar(&u.ReportPassed, "passed", false, "report unused results from functions passed to other functions")
	u.Flags.BoolVar(&u.ReportReturned, "returned", false, "report unused results from functions returned by other functions")
	u.Flags.BoolVar(&u.ReportAssigned, "assigned", false, "report unused results from functions assigned to storage")
	u.Flags.BoolVar(&u.Debug, "verbose", false, "issue debug logging")

	u.Run = u.run

	return u
}

type usage struct {
	results  uint // track explicitly used results as a bit field
	passed   bool // tracks funcs passed to other funcs
	returned bool // tracks funcs returned by other funcs
	assigned bool // tracks funcs assigned to storage
}

// callee is a target we may report. Typically a *ssa.Function or a *types.Func.
type callee interface {
	Name() string
	Pos() token.Pos
	Type() types.Type
}

func (u *unretAnalyzer) run(pass *analysis.Pass) (interface{}, error) {
	debug := log.New(log.Default().Writer(), pass.Pkg.Name()+": ", log.Lshortfile)
	if !u.Debug {
		debug.SetOutput(io.Discard)
	}
	prog := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	used := make(map[callee]usage, len(prog.SrcFuncs))

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
	funcs := make([]callee, 0, len(prog.SrcFuncs))
	typeFuncs := make(map[*types.TypeName]map[string]callee)
	for _, fun := range prog.SrcFuncs {
		funcs = append(funcs, fun)
		if recv := fun.Signature.Recv(); recv != nil {
			if typ := recvType(recv.Type()); typ != nil {
				m, ok := typeFuncs[typ]
				if !ok {
					m = make(map[string]callee)
					typeFuncs[typ] = m
				}
				m[fun.Name()] = fun
			}
		}
	}

	for _, name := range pass.Pkg.Scope().Names() {
		tn, _ := pass.Pkg.Scope().Lookup(name).(*types.TypeName)
		if tn == nil {
			continue
		}
		typ := tn.Type()
		if nam, ok := typ.(*types.Named); ok {
			typ = nam.Underlying()
		}
		if itf, ok := typ.(*types.Interface); ok {
			m, ok := typeFuncs[tn]
			if !ok {
				m = make(map[string]callee, itf.NumMethods())
				typeFuncs[tn] = m
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
	stored := make(map[ssa.Value]ssa.Value)
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
				switch instr := instr.(type) {
				case *ssa.MakeClosure:
					if len(instr.Bindings) > 0 {
						fn := instr.Fn.(*ssa.Function)
						if fn.Synthetic != "" {
							continue // ignore these; see testdata/src/g
						}
						if prev, ok := closures[fn]; ok {
							debug.Println("Repeat closure", prev, instr, pass.Fset.Position(prev.Pos()), pass.Fset.Position(instr.Pos()))
						}
						closures[fn] = instr
					}
				case *ssa.Store:
					addr := storeVal(instr.Addr)
					val := storeVal(instr.Val)
					stored[addr] = val
					// debug.Println("STORED", addr.Name(), addr, "=", val.Name(), val)
				}
			}
		}
	}

	trackAnonInterface := func(itf *types.Interface) *types.TypeName {
		anon, ok := anonTypes[itf.String()]
		if !ok {
			anon = types.NewTypeName(token.NoPos, prog.Pkg.Pkg, itf.String(), itf)
			anonTypes[itf.String()] = anon
			m := make(map[string]callee, itf.NumExplicitMethods())
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
	type index struct{ int }
	self := index{0}
	called := index{1}
	nothing := index{0}
	traversed := make(map[ssa.Value]struct{})
	var funResult func(op ssa.Value, extract func(*types.TypeName) callee) (res []callee, idx index)
	funResult = func(op ssa.Value, extract func(*types.TypeName) callee) (res []callee, idx index) {
		if _, ok := traversed[op]; ok {
			return res, nothing
		}
		traversed[op] = struct{}{}
		defer delete(traversed, op)

		switch op := op.(type) {
		case *ssa.Extract:
			if fun, _ := funResult(op.Tuple, extract); fun != nil {
				return fun, index{op.Index + 1}
			}
			return res, nothing
		case *ssa.Call:
			com := op.Common()
			if sfunc := com.StaticCallee(); sfunc != nil {
				return append(res, sfunc), called
			}
			if com.Method != nil {
				extract = func(tn *types.TypeName) callee {
					return typeFuncs[tn][com.Method.Name()]
				}
			}
			f, _ := funResult(com.Value, extract)
			return f, called
		case *ssa.MakeInterface:
			f, i := funResult(op.X, extract)
			switch t := op.Type().(type) {
			case *types.Interface:
				if ex := extract(trackAnonInterface(t)); ex != nil {
					return append(f, ex), i
				}
			case *types.Named:
				if ex := extract(t.Obj()); ex != nil {
					return append(f, ex), i
				}
			}
			return f, i
		case *ssa.Alloc:
			addr := storeVal(op)
			if val, ok := stored[addr]; ok {
				// avoid infinite recursion... not quite sure what repros it.
				delete(stored, addr)
				defer func() { stored[addr] = val }()
				return funResult(val, extract)
			} else if typ := recvType(op.Type().Underlying().(*types.Pointer).Elem()); typ != nil {
				if ex := extract(typ); ex != nil {
					return append(res, ex), self
				}
			}
			return nil, nothing
		case *ssa.Slice: // may be varargs
			return funResult(op.X, extract)
		case *ssa.Parameter:
			if typ := recvType(op.Type()); typ != nil {
				if ex := extract(typ); ex != nil {
					return append(res, ex), self
				}
			}
			if itf, ok := op.Type().(*types.Interface); ok {
				if ex := extract(trackAnonInterface(itf)); ex != nil {
					return append(res, ex), self
				}
			}
			return nil, nothing
		case *ssa.UnOp:
			return funResult(op.X, extract)
		case *ssa.Phi:
			for _, edge := range op.Edges {
				switch edge := edge.(type) {
				case *ssa.Function:
					res = append(res, edge)
				default:
					f, _ := funResult(edge, extract)
					res = append(res, f...)
				}
			}
			return res, self
		case *ssa.FreeVar:
			for i, fv := range op.Parent().FreeVars {
				if fv != op {
					continue
				}
				f, _ := funResult(closures[op.Parent()].Bindings[i], extract)
				res = append(res, f...)
			}
			return res, self

		case *ssa.MakeClosure:
			if sfn := op.Fn.(*ssa.Function); sfn.Synthetic != "" {
				if tfn, ok := sfn.Object().(*types.Func); ok {
					sig := tfn.Type().(*types.Signature)
					if fun := typeFuncs[recvType(sig.Recv().Type())][tfn.Name()]; fun != nil {
						res = append(res, fun)
					}
					return res, self
				}
			}
			if op.Fn != nil {
				return append(res, op.Fn), nothing
			}
			return res, nothing

		case *ssa.ChangeInterface:
			// debug.Println("ChangeInterface:", op.Type(), reflect.TypeOf(op.Type()), op.X, reflect.TypeOf(op.X))
			return funResult(op.X, extract)
		case *ssa.ChangeType:
			switch x := op.X.(type) {
			case *ssa.Function:
				return append(res, x), self
			}
			return funResult(op.X, extract)
		default:
			// debug.Printf("whattabout %T: %v", op, op)
			return res, nothing
		}
	}

	// apply updates records, both direct and any indirect ones.
	// inst is used only for (disabled) logging.
	apply := func(inst ssa.Instruction, fun callee, set func(*usage)) {
		debugTrackApply(debug, inst, fun, set)
		r := used[fun]
		set(&r)
		used[fun] = r
		if sfunc, ok := fun.(*ssa.Function); ok {
			if fun, ok := sfunc.Object().(*types.Func); ok {
				r := used[fun]
				set(&r)
				used[fun] = r
			}
		}
	}

	nilExtract := func(*types.TypeName) callee { return nil }
	for _, fn := range prog.SrcFuncs {
		debugDumpFuncSSA(fn, debug)
		for _, blk := range fn.Blocks {
			for _, inst := range blk.Instrs {
				switch inst := inst.(type) {
				case ssa.Value:
					// switch inst := inst.(type) {
					// case *ssa.Call, *ssa.Extract:
					funs, _ := funResult(inst, nilExtract)
					for _, fun := range funs {
						if r, ok := used[fun]; fun != nil && !ok {
							used[fun] = r
						}
					}
					// Don't consider raw extract as use of a value.
					// This lets a, _ := f() not 'use' f's second result.
					if _, ok := inst.(*ssa.Extract); ok {
						continue
					}
					// }
				case *ssa.Store:
					switch inst.Addr.(type) {
					case *ssa.FieldAddr:
						switch inst.Val.(type) {
						case *ssa.Function:
							apply(inst, inst.Val, func(r *usage) { r.assigned = true })
						}
					}
				}

				// consider other instructions referring to call/extract as function and result uses
				for _, op := range inst.Operands(nil) {
					if op == nil || *op == nil {
						continue
					}
					funs, res := funResult(*op, nilExtract)
					set := func(r *usage) { r.results |= 1 << res.int }
					if res == self {
						inner := set
						switch inst := inst.(type) {
						case *ssa.Call:
							if *op != inst.Call.Value { // callee is not passed, only args
								set = func(r *usage) { r.passed = true; inner(r) }
							}
						case *ssa.Return:
							set = func(r *usage) { r.returned = true; inner(r) }
						case *ssa.Store:
							if *op == inst.Val {
								switch inst.Addr.(type) {
								case *ssa.FieldAddr:
									set = func(r *usage) { r.assigned = true; inner(r) }
								}
							}
						}
					}
					for _, fun := range funs {
						apply(inst, fun, set)
					}
				}
			}
		}
	}

	debugDumpUses(debug, pass.Fset, used)

	for _, fn := range funcs {
		if !u.ReportExported && token.IsExported(fn.Name()) {
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
		switch {
		case !u.ReportUncalled && !ok,
			!u.ReportAssigned && r.assigned,
			!u.ReportPassed && r.passed,
			!u.ReportReturned && r.returned:
			continue
		}
		results := fn.Type().(*types.Signature).Results()
		for i, n := 0, results.Len(); i < n; i++ {
			if (1<<uint(i+1))&r.results == 0 {
				res := results.At(i)
				rnt := strings.TrimSpace(res.Name() + " " + res.Type().String())
				pass.Reportf(fn.Pos(), "%s result %d (%s) is never used", fn.Name(), i, rnt)
				// debug.Println(fn.Name(), "result", i, "unused")
			}
		}
	}

	return nil, nil
}

func storeVal(val ssa.Value) ssa.Value {
	for {
		switch v := val.(type) {
		case *ssa.IndexAddr:
			val = v.X
		default:
			return val
		}
	}
}

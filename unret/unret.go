package unret

import (
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = `unret reports returns from unexported functions that are never used.`

var reportExported, reportUncalled bool

func init() {
	Analyzer.Flags.BoolVar(&reportExported, "exported", false, "report unused returns from exported functions")
	Analyzer.Flags.BoolVar(&reportUncalled, "uncalled", false, "report unused returns from uncalled functions")
}

var Analyzer = &analysis.Analyzer{
	Name:     "unret",
	Doc:      doc,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	type returns struct {
		used uint
	}
	prog := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	f := make(map[*types.Func]returns, len(prog.SrcFuncs))

	recvType := func(t types.Type) string {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		}
		if nam, ok := t.(*types.Named); ok && nam.Obj() != nil {
			return nam.Obj().Id()
		}
		return ""
	}

	// typeFuncs stores typename->funcname->*types.Func so that methods used
	// via interface can be tracked. See testdata/b.
	funcs := make([]*types.Func, 0, len(prog.SrcFuncs))
	typeFuncs := make(map[string]map[string]*types.Func)
	for _, fun := range prog.SrcFuncs {
		if fun.Object() == nil {
			// ignore anonymous functions... (testdata/src/d)
			continue
		}
		funcs = append(funcs, fun.Object().(*types.Func))
		if recv := fun.Signature.Recv(); recv != nil {
			if typ := recvType(recv.Type()); typ != "" {
				m, ok := typeFuncs[typ]
				if !ok {
					m = make(map[string]*types.Func)
					typeFuncs[typ] = m
				}
				m[fun.Name()] = fun.Object().(*types.Func)
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
			m, ok := typeFuncs[obj.Id()]
			if !ok {
				m = make(map[string]*types.Func, itf.NumMethods())
				typeFuncs[obj.Id()] = m
			}
			for i, n := 0, itf.NumMethods(); i < n; i++ {
				fn := itf.Method(i)
				m[fn.Name()] = fn
				funcs = append(funcs, fn)
			}
		}
	}

	// for k, v := range typeFuncs {
	// 	log.Println("tf:", k, v)
	// }

	// funResult returns the *types.Func and result index that a ssa.Value
	// uses. Most of these indicate use of that result, but *ssa.Extract isn't
	// itself a use.
	var funResult func(op ssa.Value) (*types.Func, int)
	funResult = func(op ssa.Value) (*types.Func, int) {
		switch op := op.(type) {
		case *ssa.Extract:
			fun, _ := funResult(op.Tuple)
			return fun, op.Index
		case *ssa.Call:
			com := op.Common()
			if sfunc := com.StaticCallee(); sfunc != nil {
				if sfunc.Object() == nil {
					return nil, 0 // ignore anonymous functions for now
				}
				return sfunc.Object().(*types.Func), 0
			}
			if com.Method != nil {
				switch v := com.Value.(type) {
				case *ssa.MakeInterface:
					switch x := v.X.(type) {
					case *ssa.Alloc:
						if typ := recvType(x.Type()); typ != "" {
							fun := typeFuncs[typ][com.Method.Name()]
							return fun, 0
						}
					}
				case *ssa.Parameter:
					if typ := recvType(v.Type()); typ != "" {
						if tfs, ok := typeFuncs[typ]; ok {
							fun := tfs[com.Method.Name()]
							return fun, 0
						}
					}
				}
				return funResult(com.Value)
			}
			// log.Printf("%#v", op)
			return nil, 0

		case *ssa.MakeInterface:
			switch op.Type().(type) {
			case *types.Named:
				// log.Println("MakeInterface:", op.Type(), reflect.TypeOf(op.Type()), op.X, reflect.TypeOf(op.X))
				f, i := funResult(op.X)
				// log.Println("MakeInterface:", f, i)
				return f, i
			case *types.Interface:
			}
			return funResult(op.X)
		case *ssa.ChangeInterface:
			// log.Println("ChangeInterface:", op.Type(), reflect.TypeOf(op.Type()), op.X, reflect.TypeOf(op.X))
			return funResult(op.X)
		default:
			// log.Printf("whattabout %T: %v", op, op)
			return nil, 0
		}
	}

	for _, fn := range prog.SrcFuncs {
		for _, blk := range fn.Blocks {
			for _, inst := range blk.Instrs {
				if val, ok := inst.(ssa.Value); ok {
					switch val := val.(type) {
					// consider extract as function uses, but not as result uses
					case *ssa.Call:
						fun, _ := funResult(val)
						if r, ok := f[fun]; !ok {
							f[fun] = r
						}
						// keep going
					case *ssa.Extract:
						fun, _ := funResult(val)
						if r, ok := f[fun]; !ok {
							f[fun] = r
						}
						continue
					}
				}

				// consider other instructions referring to call/extract as function and result uses
				for _, op := range inst.Operands(nil) {
					if op == nil || *op == nil {
						continue
					}
					if fun, res := funResult(*op); fun != nil {
						r := f[fun]
						r.used |= 1 << res
						f[fun] = r
					}
				}
			}
		}
	}

	// log.Println(f)

	for _, fn := range funcs {
		// sb := &bytes.Buffer{}
		// ssa.WriteFunction(sb, fn)
		// io.Copy(os.Stdout, sb)
		if !reportExported && token.IsExported(fn.Name()) {
			continue
		}
		// if len(fn.Blocks[0].Instrs) > 4 {
		// 	for _, blk := range fn.Blocks {
		// 		for _, inst := range blk.Instrs {
		// 			log.Println(fn, blk, inst, reflect.TypeOf(inst))
		// 		}
		// 	}
		// }
		r, ok := f[fn]
		if !reportUncalled && !ok {
			continue
		}
		results := fn.Type().(*types.Signature).Results()
		for i, n := 0, results.Len(); i < n; i++ {
			if (1<<uint(i))&r.used == 0 {
				res := results.At(i)
				rnt := strings.TrimSpace(res.Name() + " " + res.Type().String())
				pass.Reportf(fn.Pos(), "%s result %d, %s is never used", fn.Name(), i, rnt)
			}
		}
	}

	return nil, nil
}

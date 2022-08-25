package cfappend

import "log"

// This construct generates a *ssa.Phi loop. Ensure unret doesn't crash.

func use(r []interface{}) {
	vs := r[:0]
	for _, v := range r {
		n := len(vs)
		if n == 0 || vs[n-1] != v {
			vs = append(vs, v)
		}
	}
	log.Println(vs...)
}

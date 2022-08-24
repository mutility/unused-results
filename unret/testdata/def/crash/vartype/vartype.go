package vartype

type struc struct{}

type itf interface{}

// s was once mishandled while setting up `funcs`. Its Scope entry has a .Type
// of *types.Interface (or *types.Named for one). But it is a *types.Var, not
// *types.TypeName.
var s itf = struc{}

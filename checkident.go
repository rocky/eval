package eval

import (
	"reflect"
	"go/ast"
)

func CheckIdent(ident *ast.Ident, env Env) (_ *Ident, errs []error) {
	aexpr := &Ident{Ident: ident}
	switch aexpr.Name {
	case "nil":
		aexpr.constValue = ConstValueOf(UntypedNil{})
		aexpr.knownType = []reflect.Type{ConstNil}

	case "true":
		aexpr.constValue = ConstValueOf(true)
		aexpr.knownType = []reflect.Type{ConstBool}

	case "false":
		aexpr.constValue = ConstValueOf(false)
		aexpr.knownType = []reflect.Type{ConstBool}
	default:
		for searchEnv := env; searchEnv != nil; searchEnv = searchEnv.PopScope() {
			if v := searchEnv.Var(aexpr.Name); v.IsValid() {
				aexpr.knownType = knownType{v.Elem().Type()}
				aexpr.source = EnvVar
				return aexpr, errs
			} else if v := searchEnv.Func(aexpr.Name); v.IsValid() {
				aexpr.knownType = knownType{v.Type()}
				aexpr.source = EnvFunc
				return aexpr, errs
			} else if v := searchEnv.Const(aexpr.Name); v.IsValid() {
				if n, ok := v.Interface().(*ConstNumber); ok {
					aexpr.knownType = knownType{n.Type}
				} else {
					aexpr.knownType = knownType{v.Type()}
				}
				aexpr.constValue = constValue(v)
				aexpr.source = EnvConst
				return aexpr, errs
			}
		}
		return aexpr, append(errs, ErrUndefined{aexpr})
        }
	return aexpr, errs
}

package eval

import (
	"reflect"
	"go/ast"
)

type CheckIdentFn func(ident *ast.Ident, env Env) (_ *Ident, errs []error)
var checkIdentFn CheckIdentFn

func SetCheckIdent(fn CheckIdentFn) {
	checkIdentFn = fn
}


type EvalIdentFn func(ident *Ident, env Env) (reflect.Value, error)
var evalIdentFn EvalIdentFn

func SetEvalIdent(fn EvalIdentFn) {
	evalIdentFn = fn
}

type CheckSelectorExprFn func(selector *ast.SelectorExpr, env Env) (*SelectorExpr, []error)
var checkSelectorExprFn CheckSelectorExprFn

func SetCheckSelectorExpr(fn CheckSelectorExprFn) {
	checkSelectorExprFn = fn
}

type EvalSelectorExprFn func(selector *SelectorExpr, env Env) (reflect.Value, error)

var evalSelectorExprFn EvalSelectorExprFn

func SetEvalSelectorExpr(fn EvalSelectorExprFn) {
	evalSelectorExprFn = fn
}

func init() {
	SetCheckIdent(checkIdent)
	SetEvalIdent(evalIdent)
	SetCheckSelectorExpr(checkSelectorExpr)
	SetEvalSelectorExpr(evalSelectorExpr)
}

package eval

// Tests replacing the default identifier selection lookup value mechanism with
// our own custom versions.

import (
	"go/ast"
	"reflect"
	"testing"
)

// Here's our custom ident lookup.
func myEvalIdent(ident *Ident, env Env) (reflect.Value, error) {
	name := ident.Name
	if name == "nil" {
		return EvalNil, nil
	} else if name == "a" {
		val := reflect.ValueOf(5)
		return val, nil
	} else if name[0] == 'v' {
		val := reflect.ValueOf(5)
		return val, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf("constant")
		return val, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf(true)
		return val, nil
	} else {
		val := reflect.ValueOf('x')
		return val, nil
	}
}


// Here's our custom ident type check
func myCheckIdent(ident *ast.Ident, env Env) (_ *Ident, errs []error) {
	aexpr := &Ident{Ident: ident}
	name := aexpr.Name
	switch name {
	case "nil":
		aexpr.SetConstValue(ConstValueOf(UntypedNil{}))
		aexpr.SetKnownType([]reflect.Type{ConstNil})
	case "true":
		aexpr.SetConstValue(ConstValueOf(true))
		aexpr.SetKnownType([]reflect.Type{ConstBool})
	case "false":
		aexpr.SetConstValue(ConstValueOf(false))
		aexpr.SetKnownType([]reflect.Type{ConstBool})
	case "a":
		val := reflect.ValueOf(5)
		knowntype := knownType{val.Type()}
		aexpr.SetKnownType(knowntype)
		aexpr.SetSource(EnvVar)
	default:
		if name[0] == 'v' {
			aexpr.SetConstValue(ConstValueOf(true))
			aexpr.SetKnownType([]reflect.Type{ConstBool})
		} else if name[0] == 'c' {
			aexpr.knownType = knownType{stringType}
			aexpr.source = EnvConst
		} else if name == "bogus" {
			aexpr.source = EnvConst
		} else {
			aexpr.knownType = knownType{f32}
		aexpr.source = EnvVar
		}
	}
	return aexpr, errs
}


// Here's our custom selector lookup.
func myEvalSelectorExpr(selector *SelectorExpr, env Env) (
	reflect.Value, error) {
	val := reflect.ValueOf("bogus")
	return val, nil
}

// Here's our custom selector type check
func myCheckSelectorExpr(selector *ast.SelectorExpr, env Env) (*SelectorExpr, []error) {
	aexpr := &SelectorExpr{SelectorExpr: selector}
	sel, errs := myCheckIdent(selector.Sel, env)
	aexpr.constValue = sel.constValue
	aexpr.knownType = sel.knownType
	return aexpr, errs
}

func TestReplaceEval(t *testing.T) {
	defer SetEvalIdent(EvalIdent)
	defer SetCheckIdent(CheckIdent)
	env := MakeSimpleEnv()
	SetCheckIdent(myCheckIdent)
	expectResult(t, "true || false", env, true)
	expectResult(t, "true && false", env, false)

	SetEvalIdent(myEvalIdent)
	expectResult(t, "fdafdsa", env, 'x')
	expectResult(t, "c + \" value\"", env, "constant value")
	expectResult(t, "a+5", env, 10)

}

func TestReplaceSelectorLookup(t *testing.T) {
	defer SetCheckSelectorExpr(CheckSelectorExpr)
	defer SetEvalSelectorExpr(EvalSelectorExpr)
	SetCheckSelectorExpr(myCheckSelectorExpr)
	SetEvalSelectorExpr(myEvalSelectorExpr)
	env  := MakeSimpleEnv()
	pkg := MakeSimpleEnv()
	env.Pkgs["bogusPackage"] = pkg
	expectResult(t, "bogusPackage.something", env, "bogus")
}

// Demos replacing the default identifier lookup value mechanism with
// our own custom version.

package main

import (
	"go/ast"
	"reflect"
	"fmt"
	"github.com/rocky/eval"
)

// Here's our custom ident lookup.
func EvalIdent(ident *eval.Ident, env eval.Env) (reflect.Value, error) {
	println("Evaldent called")
	name := ident.Name
	if name == "nil" {
		return eval.EvalNil, nil
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
func CheckIdent(ident *ast.Ident, env eval.Env) (_ *eval.Ident, errs []error) {
	println("CheckIdent called")
	aexpr := &eval.Ident{Ident: ident}
	name := aexpr.Name
	switch name {
	case "nil":
		aexpr.SetConstValue(eval.ConstValueOf(eval.UntypedNil{}))
		aexpr.SetKnownType([]reflect.Type{eval.ConstNil})
		return aexpr, errs
	case "true":
		aexpr.SetConstValue(eval.ConstValueOf(true))
		aexpr.SetKnownType([]reflect.Type{eval.ConstBool})
		return aexpr, errs
	case "false":
		aexpr.SetConstValue(eval.ConstValueOf(false))
		aexpr.SetKnownType([]reflect.Type{eval.ConstBool})
	default:
		errs = append(errs, eval.ErrUndefined{aexpr})
		return aexpr, errs
	}
	return aexpr, errs
}

var evalEnv eval.Env = eval.MakeSimpleEnv()

func EvalExpr(expr string) ([]reflect.Value, error) {
	results, panik, compileErrs := eval.EvalEnv(expr, evalEnv)
	if compileErrs != nil {
		println("compileErr != nil", )
		for _, err := range(compileErrs) {
			fmt.Printf("+++ %T\n", err)
			fmt.Printf("+++2 %v\n", err)
		}
	} else if panik != nil {
		println("panic != nil")
		for _, err := range(compileErrs) {
			fmt.Println(err.Error())
		}
	} else {
		println("ok")
		return results, nil
	}
	return nil, nil
}

func main() {
	if results, err := EvalExpr("true"); err == nil {
		fmt.Printf("%v\n", results[0].Interface())
	}
	eval.SetCheckIdent(CheckIdent)
	// eval.SetEvalIdent(EvalIdent)
	if results, err := EvalExpr("true"); err == nil {
		fmt.Printf("%v\n", results[0].Interface())
	}
	if results, err := EvalExpr("true || false"); err == nil {
	 	println("true || false")
	 	fmt.Printf("%v\n", results[0].Interface())
	}
	if results, err := EvalExpr("true && false"); err == nil {
		println("true && false")
		fmt.Printf("%v\n", results[0].Interface())
	}

}

// This shows how to replace the default identifier lookup and
// id selector lookup with custom routines.
//
// This is used in the gub debugger where
// the environment structure and interpreter values
// are different than what eval uses, but we still want
// to use eval for its ability to parse, type check,
// and run expressions.

package main

import (
	"errors"
	"go/ast"
	"reflect"
	"fmt"
	"os"
	"github.com/rocky/eval"
)

type knownType []reflect.Type

func makeBogusEnv() *eval.SimpleEnv {

	// A copule of things from the fmt package.

	// A stripped down package environment.  See
	// http://github.com/rocky/go-fish and repl_imports.go for a more
	// complete environment.
	pkgs := map[string] eval.Env {
			"fmt": &eval.SimpleEnv {
				Vars:   make(map[string] reflect.Value),
				Consts: make(map[string] reflect.Value),
				Funcs:  make(map[string] reflect.Value),
				Types : make(map[string] reflect.Type),
				Pkgs:   nil,
			}, "os": &eval.SimpleEnv {
				Vars:   map[string] reflect.Value {
					"Stdout": reflect.ValueOf(&os.Stdout),
					"Args"  : reflect.ValueOf(&os.Args)},
				Consts: make(map[string] reflect.Value),
				Funcs:  make(map[string] reflect.Value),
				Types:  make(map[string] reflect.Type),
				Pkgs:   nil,
			},
		}

	mainEnv := eval.MakeSimpleEnv()
	mainEnv.Pkgs = pkgs

	a := 5

	mainEnv.Vars["a"]    = reflect.ValueOf(&a)

	return mainEnv
}

func pkgEvalIdent(ident *eval.Ident, env eval.Env) (reflect.Value, error) {
	if ident.IsConst() {
		return ident.Const(), nil
	}

	name := ident.Name
	switch ident.Source() {
	case eval.EnvVar:
		for searchEnv := env; searchEnv != nil; searchEnv = searchEnv.PopScope() {
			if v := searchEnv.Var(name); v.IsValid() {
				return v.Elem(), nil
			}
		}
	case eval.EnvFunc:
		for searchEnv := env; searchEnv != nil; searchEnv = searchEnv.PopScope() {
			if v := searchEnv.Func(name); v.IsValid() {
				return v, nil
			}
		}
	}
	return reflect.Value{}, errors.New("Something went wrong")
}

// Here's our custom ident lookup.
func EvalIdent(ident *eval.Ident, env eval.Env) (reflect.Value, error) {
	name := ident.Name
	fmt.Printf("EvalIdent %s called\n", name)
	if name == "nil" {
		return eval.EvalNil, nil
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
	}
	return eval.EvalIdent(ident, env)

}

// Here's our custom ident type check
func CheckIdent(ident *ast.Ident, env eval.Env) (_ *eval.Ident, errs []error) {
	aexpr := &eval.Ident{Ident: ident}
	name := aexpr.Name
	fmt.Printf("CheckIdent %s called\n", name)
	switch name {
	case "nil", "true", "false":
		return eval.CheckIdent(ident, env)
	case "a":
		val := reflect.ValueOf(5)
		knowntype := knownType{val.Type()}
		aexpr.SetKnownType(knowntype)
		aexpr.SetSource(eval.EnvVar)
	default:
		return eval.CheckIdent(ident, env)
	}
	return aexpr, errs
}

// Here's our custom selector lookup.
func EvalSelectorExpr(selector *eval.SelectorExpr, env eval.Env) (reflect.Value, error) {
	println("custom EvalSelectorExpr called")

	if pkgName := selector.PkgName(); pkgName != "" {
		if vs, err := pkgEvalIdent(selector.Sel, env.Pkg(pkgName)); err != nil {
			return EvalIdent(selector.Sel, env.Pkg(selector.PkgName()))
		} else {
			return vs, err
		}
	}

	vs, err := eval.EvalExpr(selector.X, env)
	if err != nil {
		return reflect.Value{}, err
	}

	v := vs[0]
	t := v.Type()
	if selector.Field() != nil {
		if t.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.FieldByIndex(selector.Field()), nil
	}

	if selector.IsPtrReceiver() {
		v = v.Addr()
	}
	return v.Method(selector.Method()), nil
}

var evalEnv eval.Env = makeBogusEnv()

func EvalExpr(expr string) ([]reflect.Value, error) {
	results, panik, compileErrs := eval.EvalEnv(expr, evalEnv)
	if compileErrs != nil {
		fmt.Println("compile errors:")
		for _, err := range(compileErrs) {
			fmt.Println(err.Error())
		}
	} else if panik != nil {
		fmt.Printf("Evaluation panic: %s\n", panik.Error())
	} else {
		return results, nil
	}
	return nil, nil
}

func main() {
	// if results, err := EvalExpr("a"); err == nil {
	// 	fmt.Printf("%v\n", results[0].Interface())
	// } else {
	// 	println("Can't eval 'a'")
	// }
	eval.SetCheckIdent(CheckIdent)
	eval.SetEvalIdent(EvalIdent)
	eval.SetEvalSelectorExpr(EvalSelectorExpr)
	// if results, err := EvalExpr("a+5"); err == nil {
	// 	fmt.Printf("%v\n", results[0].Interface())
	// }

	if results, err := EvalExpr("os.Args"); err == nil {
		fmt.Printf("%v\n", results[0].Interface())
	}

	// if results, err := EvalExpr("true"); err == nil {
	// 	fmt.Printf("%v\n", results[0].Interface())
	// }
	// if results, err := EvalExpr("true || false"); err == nil {
	//  	println("true || false")
	//  	fmt.Printf("%v\n", results[0].Interface())
	// }
	// if results, err := EvalExpr("true && false"); err == nil {
	// 	println("true && false")
	// 	fmt.Printf("%v\n", results[0].Interface())
	// }

}

package go_parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

func parseAndEval(exp string) (int, error) {
	tree, err := parser.ParseExpr(exp)
	if err != nil {
		return 0, err
	}
	return eval(tree)
}

func eval(tree ast.Expr) (int, error) {
	switch n := tree.(type) {
	case *ast.BasicLit:
		if n.Kind != token.INT {
			return unsup(n.Kind)
		}
		i, _ := strconv.Atoi(n.Value)
		return i, nil
	case *ast.BinaryExpr:
		switch n.Op {
		case token.ADD, token.SUB, token.MUL, token.QUO:
		default:
			return unsup(n.Op)
		}
		x, err := eval(n.X)
		if err != nil {
			return 0, err
		}
		y, err := eval(n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			return x / y, nil
		}
	case *ast.ParenExpr:
		return eval(n.X)
	}
	return unsup(reflect.TypeOf(tree))
}

func unsup(i interface{}) (int, error) {
	return 0, errors.New(fmt.Sprintf("%v unsupported", i))
}

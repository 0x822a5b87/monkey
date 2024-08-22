package evaluator

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/object"
	"fmt"
	"reflect"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expr)
	case *ast.IntegerLiteral:
		return evalIntegralLiteral(node)
	default:
		panic(fmt.Errorf("error node type for [%s]", reflect.TypeOf(node).String()))
	}
}

func evalIntegralLiteral(integerLiteral *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: integerLiteral.Value}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		// TODO a sequence of stmt should be evaluate
		result = Eval(stmt)
	}
	return result
}

package evaluator

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/common"
	"0x822a5b87/monkey/object"
	"0x822a5b87/monkey/token"
	"fmt"
	"reflect"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expr)
	case *ast.BooleanExpression:
		return evalBooleanLiteral(node)
	case *ast.IntegerLiteral:
		return evalIntegralLiteral(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	default:
		panic(fmt.Errorf("error node type for [%s]", reflect.TypeOf(node).String()))
	}
}

func evalIntegralLiteral(integerLiteral *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: integerLiteral.Value}
}

func evalBooleanLiteral(booleanExpression *ast.BooleanExpression) object.Object {
	return nativeBoolean(booleanExpression.Value)
}

func evalInfixExpression(infix *ast.InfixExpression) object.Object {
	lhsObj := Eval(infix.Lhs)
	rhsObj := Eval(infix.Rhs)
	switch infix.Operator {
	case string(token.SUB):
		return evalSubtract(lhsObj, rhsObj)
	case string(token.ASTERISK):
		return evalMultiply(lhsObj, rhsObj)
	case string(token.GT):
		return evalGreaterThan(lhsObj, rhsObj)
	case string(token.LT):
		return evalLessThan(lhsObj, rhsObj)
	case string(token.EQ):
		return evalEqual(lhsObj, rhsObj)
	case string(token.NotEq):
		return evalNotEqual(lhsObj, rhsObj)
	}

	// TODO support more infix expression
	return object.NativeNull
}

func evalNotEqual(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Equatable)
	r := rhsObj.(object.Equatable)
	return l.NotEqual(r)
}

func evalEqual(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Equatable)
	r := rhsObj.(object.Equatable)
	return l.Equal(r)
}

func evalLessThan(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Comparable)
	r := rhsObj.(object.Comparable)
	return l.LessThan(r)
}

func evalGreaterThan(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Comparable)
	r := rhsObj.(object.Comparable)
	return l.GreaterThan(r)
}

func evalSubtract(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Subtract)
	r := rhsObj.(object.Subtract)
	return l.Sub(r)
}

func evalMultiply(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Multiply)
	r := rhsObj.(object.Multiply)
	return l.Mul(r)
}

func evalInfixExpressionIntegerLiteral(operator token.TokenType, lhsIntegerObj, rhsIntegerObj *object.Integer) object.Object {
	switch operator {
	case token.SUB:
		return &object.Integer{Value: lhsIntegerObj.Value - rhsIntegerObj.Value}
	case token.ASTERISK:
		return &object.Integer{Value: lhsIntegerObj.Value * rhsIntegerObj.Value}
	default:
		return object.NativeNull
	}
}

func evalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	switch prefix.Operator {
	case string(token.BANG):
		return evalBangOfPrefixExpression(prefix.Right)
	case string(token.SUB):
		return evalMinusOfPrefixExpression(prefix.Right)
	default:
		panic(common.ErrUnknownToken)
	}
}

func evalMinusOfPrefixExpression(rightExpr ast.Expression) object.Object {
	right := Eval(rightExpr)
	if right.Type() != object.ObjInteger {
		// TODO think of it, return NativeNull or panic
		return object.NativeNull
	}
	integer := right.(*object.Integer)
	return &object.Integer{Value: -integer.Value}
}

func evalBangOfPrefixExpression(rightExpr ast.Expression) object.Object {
	right := Eval(rightExpr)
	switch right {
	case object.NativeFalse:
		return object.NativeTrue
	case object.NativeNull:
		return object.NativeTrue
	case object.NativeTrue:
		return object.NativeFalse
	default:
		return object.NativeFalse
	}
}

func nativeBoolean(input bool) object.Object {
	if input {
		return object.NativeTrue
	} else {
		return object.NativeFalse
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		// TODO a sequence of stmt should be evaluate
		result = Eval(stmt)
	}
	return result
}

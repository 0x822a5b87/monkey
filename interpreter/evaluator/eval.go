package evaluator

import (
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/token"
	"fmt"
	"reflect"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env, false)
	case *ast.ExpressionStatement:
		return Eval(node.Expr, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env, true)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.BooleanExpression:
		return evalBooleanLiteral(node)
	case *ast.IntegerLiteral:
		return evalIntegralLiteral(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FnLiteral:
		return evalFnLiteral(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.HashExpression:
		return evalHash(node, env)
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

func evalInfixExpression(infix *ast.InfixExpression, env *object.Environment) object.Object {
	lhsObj := Eval(infix.Lhs, env)
	rhsObj := Eval(infix.Rhs, env)

	err := InfixExpressionTypeCheck(infix.Operator, lhsObj, rhsObj)
	if err != nil {
		return err
	}

	switch infix.Operator {
	case string(token.PLUS):
		return evalAdd(lhsObj, rhsObj)
	case string(token.SUB):
		return evalSubtract(lhsObj, rhsObj)
	case string(token.ASTERISK):
		return evalMultiply(lhsObj, rhsObj)
	case string(token.SLASH):
		return evalDivide(lhsObj, rhsObj)
	case string(token.GT):
		return evalGreaterThan(lhsObj, rhsObj)
	case string(token.LT):
		return evalLessThan(lhsObj, rhsObj)
	case string(token.EQ):
		return evalEqual(lhsObj, rhsObj)
	case string(token.NotEq):
		return evalNotEqual(lhsObj, rhsObj)
	case string(token.LBRACKET):
		return evalNotEqual(lhsObj, rhsObj)
	}

	// TODO support more infix expression
	return object.NativeNull
}

//func evalIndex(array, index object.Object) object.Object {
//
//}

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

func evalAdd(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Add)
	r := rhsObj.(object.Add)
	return l.Add(r)
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

func evalDivide(lhsObj, rhsObj object.Object) object.Object {
	l := lhsObj.(object.Divide)
	r := rhsObj.(object.Divide)
	return l.Divide(r)
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

func evalPrefixExpression(prefix *ast.PrefixExpression, env *object.Environment) object.Object {
	rhs := Eval(prefix.Right, env)
	err := prefixExpressionTypeCheck(prefix.Operator, rhs)
	if err != nil {
		return err
	}

	switch prefix.Operator {
	case string(token.BANG):
		return evalBangOfPrefixExpression(prefix.Right, env)
	case string(token.SUB):
		return evalMinusOfPrefixExpression(prefix.Right, env)
	default:
		panic(common.ErrUnknownToken)
	}
}

func evalMinusOfPrefixExpression(rightExpr ast.Expression, env *object.Environment) object.Object {
	right := Eval(rightExpr, env)
	integer := right.(*object.Integer)
	return &object.Integer{Value: -integer.Value}
}

func evalBangOfPrefixExpression(rightExpr ast.Expression, env *object.Environment) object.Object {
	right := Eval(rightExpr, env)
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

func evalLetStatement(letStatement *ast.LetStatement, env *object.Environment) object.Object {
	obj := Eval(letStatement.Value, env)
	if obj.Type() == object.ObjError {
		return obj
	}
	env.Set(letStatement.Name.Value, obj)
	return obj
}

func evalCallExpression(call *ast.CallExpression, env *object.Environment) object.Object {
	fnOrBuiltIn := Eval(call.Fn, env)
	switch fnValue := fnOrBuiltIn.(type) {
	case *object.Fn:
		return evalFn(call, fnValue, env)
	case *object.BuiltIn:
		return evalBuiltIn(call, fnValue, env)
	default:
		return object.NativeNull
	}
}

func evalStringLiteral(stringLiteral *ast.StringLiteral) object.Object {
	return &object.StringObj{
		Value: stringLiteral.Literal,
	}
}

func evalArrayLiteral(al *ast.ArrayLiteral, environment *object.Environment) object.Object {
	array := &object.Array{Elements: make([]object.Object, 0)}
	for _, element := range al.Elements {
		obj := Eval(element, environment)
		if obj.Type() == object.ObjError {
			return obj
		}
		array.Elements = append(array.Elements, obj)
	}
	return array
}

func evalIndexExpression(ie *ast.IndexExpression, environment *object.Environment) object.Object {
	lhs := Eval(ie.Lhs, environment)
	if lhs.Type() == object.ObjError {
		return lhs
	}

	index := Eval(ie.Index, environment)
	if index.Type() == object.ObjError {
		return index
	}

	array, ok := lhs.(object.Index)
	if !ok {
		return newError("%snot an index expression : %s", unknownOperatorErrStr, lhs.Type())
	}

	return array.Index(index)
}

func evalHash(expr *ast.HashExpression, environment *object.Environment) object.Object {
	hash := &object.Hash{Pairs: map[object.HashKey]*object.HashPair{}}
	for k, v := range expr.Pairs {
		key := Eval(k, environment)
		if key.Type() == object.ObjError {
			return key
		}
		value := Eval(v, environment)
		if value.Type() == object.ObjError {
			return value
		}

		hashable, ok := key.(object.Hashable)
		if !ok {
			return newError("%s type = [%s]", hashableNotImplementError, key.Type())
		}
		hash.Pairs[hashable.HashKey()] = &object.HashPair{
			Key:   key,
			Value: value,
		}
	}
	return hash
}

func evalFnLiteral(fnLiteral *ast.FnLiteral, env *object.Environment) *object.Fn {
	return &object.Fn{
		Params: fnLiteral.Parameters,
		Body:   fnLiteral.Body,
		Env:    env,
	}
}

func evalIdentifier(identifier *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(identifier.Value)
	if !ok {
		return newError("%s %s", identifierNotFoundErrStr, identifier.Value)
	}
	return value
}

func evalReturnStatement(returnStmt *ast.ReturnStatement, env *object.Environment) object.Object {
	return &object.Return{
		Object: Eval(returnStmt.ReturnValue, env),
	}
}

func evalStatements(stmts []ast.Statement, env *object.Environment, wrapReturn bool) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result == nil {
			continue
		}

		if result.Type() == object.ObjReturn {
			if wrapReturn {
				return result
			} else {
				r := result.(*object.Return)
				return r.Object
			}
		}

		if result.Type() == object.ObjError {
			return result
		}
	}
	return result
}

func evalIfExpression(ifStmt *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifStmt.Condition, env)
	if isTruthyObject(condition) {
		return Eval(ifStmt.Consequence, env)
	}

	if ifStmt.Alternative != nil {
		return Eval(ifStmt.Alternative, env)
	} else {
		return object.NativeNull
	}
}

//func isTruthyObject(o object.Object) bool {
//	if o.Type() == object.ObjNull {
//		return false
//	}
//
//	// an object of type Boolean with a value of false means it is not truthy
//	b, ok := o.(*object.Boolean)
//	if ok {
//		return b.Value
//	}
//
//	return true
//}

func isTruthyObject(o object.Object) bool {
	switch o {
	case object.NativeNull:
		fallthrough
	case object.NativeFalse:
		return false
	default:
		return true
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalBuiltIn(call *ast.CallExpression, builtIn *object.BuiltIn, env *object.Environment) object.Object {
	args := make([]object.Object, 0)
	for _, argument := range call.Arguments {
		argValue := Eval(argument, env)
		args = append(args, argValue)
	}
	return builtIn.BuiltInFn(args...)
}

//func evalArguments(call *ast.CallExpression, env *object.Environment) object.Object {
//
//}

func evalFn(call *ast.CallExpression, fn *object.Fn, callEnv *object.Environment) object.Object {
	if len(call.Arguments) != len(fn.Params) {
		return newError("%s expected [%d], got [%d]", paramsNumberMismatchErrStr, len(call.Arguments), len(fn.Params))
	}

	// there are two distinct environments associated with a function
	// 1. call environment : the environment when the function is called;
	// 2. function define environment : the environment when the function is defined.
	// when a function is called, there should have an environment chain:
	// inside function env -> args -> call env -> ...
	// and a function call contains its own environment and this environment is initiated when a FnLiteral is evaluated which is the fn.Env
	// the fn.Env contains environment: call env -> ...
	// so what we should do is initiate the args, and use fn.Env as its parent env.

	// env for arguments
	argumentsEnv := object.NewEnvironment(fn.Env)
	for i, arg := range call.Arguments {
		// TODO add scope for block
		value := Eval(arg, callEnv)
		if value.Type() == object.ObjError {
			return value
		}
		// bind argument value to params
		argumentsEnv.Set(fn.Params[i].String(), value)
	}
	fnEvalResult := Eval(fn.Body, argumentsEnv)
	return unwrapReturnValue(fnEvalResult)
}

//func getFnOrBuiltIn(call *ast.CallExpression, env *object.Environment) (object.Object, *object.Error) {
//	fnName, ok := call.Fn.(*ast.Identifier)
//	if ok {
//		fn := Eval(fnName, env)
//		switch function := fn.(type) {
//		case *object.BuiltIn:
//			return function, nil
//		case *object.Fn:
//			return function, nil
//		case *object.Error:
//			return nil, function
//		}
//		return nil, newError("%s from %s to %s", typeMismatchErrStr, fn.Type(), object.ObjFunction)
//	}
//	// maybe a closure
//	fnLiteral, ok := call.Fn.(*ast.FnLiteral)
//	if ok {
//		fn := Eval(fnLiteral, env)
//		return fn, nil
//	}
//
//	panic("unknown call expression")
//}

func unwrapReturnValue(obj object.Object) object.Object {
	returnObj, ok := obj.(*object.Return)
	if ok {
		return returnObj.Object
	}
	return obj
}

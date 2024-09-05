package evaluator

import (
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/token"
	"bytes"
	"fmt"
	"reflect"
)

const (
	typeMismatchErrStr         = "type mismatch:"
	unknownOperatorErrStr      = "unknown operator:"
	identifierNotFoundErrStr   = "identifier not found:"
	paramsNumberMismatchErrStr = "number of parameters mismatch:"
	hashableNotImplementError  = "hashable not implement:"
)

var infixOperatorTypes map[string]any
var prefixOperatorTypes map[string]any

func infixExpressionTypeCheck(operator string, lhs, rhs object.Object) *object.Error {

	if lhs.Type() == object.ObjError {
		return lhs.(*object.Error)
	}

	if rhs.Type() == object.ObjError {
		return rhs.(*object.Error)
	}

	interfaceValue := infixOperatorTypes[operator]

	typeMismatchErr := checkForTypeMismatchOperator(operator, lhs, rhs)
	if typeMismatchErr != nil {
		return typeMismatchErr
	}

	unknownOperatorErr := check4UnknownOperator(operator, interfaceValue, lhs, rhs)
	if unknownOperatorErr != nil {
		return unknownOperatorErr
	}
	return nil
}

func prefixExpressionTypeCheck(operator string, operand object.Object) *object.Error {
	if operand.Type() == object.ObjError {
		return operand.(*object.Error)
	}

	typeMismatchErr := checkForTypeMismatchOperator(operator, operand)
	if typeMismatchErr != nil {
		return typeMismatchErr
	}

	interfaceValue := prefixOperatorTypes[operator]
	unknownOperatorErr := check4UnknownOperator(operator, interfaceValue, operand)
	if unknownOperatorErr != nil {
		return unknownOperatorErr
	}
	return nil
}

func checkForTypeMismatchOperator(operator string, objects ...object.Object) *object.Error {
	if len(objects) == 0 {
		return nil
	}
	basicType := reflect.TypeOf(objects[0])
	for _, o := range objects {
		curType := reflect.TypeOf(o)
		if basicType != curType {
			return newTypeMismatchError(operator, objects...)
		}
	}
	return nil
}

// check4UnknownOperator 检查objects是否都实现了接口
func check4UnknownOperator(operator string, i any, objects ...object.Object) *object.Error {
	interfaceType := reflect.TypeOf(i).Elem()
	for _, o := range objects {
		objType := reflect.TypeOf(o)
		if interfaceType.Kind() != reflect.Interface {
			panic("The first argument must be an interface type")
		}
		implements := objType.Implements(interfaceType)
		if !implements {
			return newUnknownOperatorError(operator, objects...)
		}
	}
	return nil
}

func newTypeMismatchError(operator string, objects ...object.Object) *object.Error {
	buffer := bytes.Buffer{}
	buffer.WriteString(typeMismatchErrStr)
	for i, o := range objects {
		buffer.WriteString(fmt.Sprintf(" %s", o.Type()))
		if i != len(objects)-1 {
			buffer.WriteString(" ")
			buffer.WriteString(operator)
		}
	}
	return newError(buffer.String())
}

func newUnknownOperatorError(operator string, objects ...object.Object) *object.Error {
	if len(objects) == 2 {
		return newInfixUnknownOperatorError(operator, objects...)
	} else {
		return newPrefixUnknownOperatorError(operator, objects[0])
	}
}

func newInfixUnknownOperatorError(operator string, objects ...object.Object) *object.Error {
	buffer := bytes.Buffer{}
	buffer.WriteString(unknownOperatorErrStr)
	for i, o := range objects {
		buffer.WriteString(fmt.Sprintf(" %s", o.Type()))
		if i != len(objects)-1 {
			buffer.WriteString(" ")
			buffer.WriteString(operator)
		}
	}
	return newError(buffer.String())
}

func newPrefixUnknownOperatorError(operator string, o object.Object) *object.Error {
	buffer := bytes.Buffer{}
	buffer.WriteString(unknownOperatorErrStr)
	buffer.WriteString(fmt.Sprintf(" %s%s", operator, o.Type()))
	return newError(buffer.String())
}

func init() {
	infixOperatorTypes = make(map[string]any)
	infixOperatorTypes[string(token.PLUS)] = (*object.Add)(nil)
	infixOperatorTypes[string(token.SUB)] = (*object.Subtract)(nil)
	infixOperatorTypes[string(token.ASTERISK)] = (*object.Multiply)(nil)
	infixOperatorTypes[string(token.SLASH)] = (*object.Divide)(nil)
	infixOperatorTypes[string(token.GT)] = (*object.Comparable)(nil)
	infixOperatorTypes[string(token.LT)] = (*object.Comparable)(nil)
	infixOperatorTypes[string(token.EQ)] = (*object.Equatable)(nil)
	infixOperatorTypes[string(token.NotEq)] = (*object.Equatable)(nil)

	prefixOperatorTypes = make(map[string]any)
	prefixOperatorTypes[string(token.SUB)] = (*object.Negative)(nil)
}

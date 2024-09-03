### 4 - Extending the Interpreter

## 4.1 - Data Types & Functions

We will add new token types, modify the lexer, extend the parser and finlally add support for the data types to our evaluator and the object system.Even better is that the data type we're going to add are already present is Go.That means that we only need to make them available in Monky.

In addition to that we're also going to make the interpreter much more powerful by adding some new functions.These new ones, called `built-in functions`.

## 4.2 - Strings

### Supporting Strings in our Lexer

The first thing we have to do is add support for strings literals to our lexer. The basic structure of strings is this:

```mermaid
---
title: strings literal
---
flowchart LR

left_quote['""']:::operator --> Strings["sequence of characters"] --> right_quote['""']:::operator

classDef operator fill:#f9f,stroke:#333,stroke-width:4px;
```

the following code show the parsing function:

```go
func (l *Lexer) readString() string {
	// skip left quote
	l.readChar()

	buffer := bytes.Buffer{}

	var end = false
	// actually, we should parse the string with a state machine instead of peek char
	for !end {
		switch l.curCh() {
		case '\\':
			// in this case, whatever the next character is, we simply consume it as a basic char
			l.readChar()
			buffer.WriteByte(l.curCh())
			l.readChar()
		case '"':
			end = true
		default:
			buffer.WriteByte(l.curCh())
			l.readChar()
		}
	}

	return buffer.String()
}
```

### 4.3 - Built-in Functions

```go
func evalCallExpression(call *ast.CallExpression, env *object.Environment) object.Object {
	fnOrBuiltIn, err := getFnOrBuiltIn(call, env)
	if err != nil {
		return err
	}

	switch fnValue := fnOrBuiltIn.(type) {
	case *object.Fn:
		return evalFn(call, fnValue)
	case *object.BuiltIn:
		return evalBuiltIn(call, fnValue, env)
	default:
		return object.NativeNull
	}
}

func getFnOrBuiltIn(call *ast.CallExpression, env *object.Environment) (object.Object, *object.Error) {
	fnName, ok := call.Fn.(*ast.Identifier)
	if ok {
		fn := Eval(fnName, env)
		switch function := fn.(type) {
		case *object.BuiltIn:
			return function, nil
		case *object.Fn:
			return function, nil
		}
		return nil, newError("%s from %s to %s", typeMismatchErrStr, fn.Type(), object.ObjFunction)
	}
	// maybe a closure
	fnLiteral, ok := call.Fn.(*ast.FnLiteral)
	if ok {
		fn := Eval(fnLiteral, env)
		return fn, nil
	}

	panic("unknown call expression")
}

```

### 4.4 - Array

The data type we're going to add to our monkey interpreter in this section is the array.

```js
let monkeyArray = ["hello", "world", 28, fn(x) { x * x };];
```

In this section we'll alose add support for arrays to our newly added `len` function and also add a few more built-in functions that work with arrays.

```js
let monekyArray = ["one", "two", "three"];
len(monkeyArray);
```

The basis for our implementation of the monkey array in our interpreter wil be a Go slice of type `[]object.Object`.

### 4.5 - Hashes

A hash is what's sometimes called `hash`, `map`, `hash map` or `dictionary` in programming languages.It maps key to values.
























































# monkey

## acknowledge

This is an experimental project under the guidance of [Writing An Interpreter In Go](https://interpreterbook.com/) and [Writing A Compiler In Go](https://compilerbook.com/).

In this repository, I am going to build my own `lexer`,  `parser`, `AST` , `evaluator`, `compiler`, `Intermediate Representation`, `Virtual Machine` .

Pay my utmost tribute to the author [Thorsten Ball](https://thorstenball.com/), an extraordinarily excellent project you have created!

## install

```bash
go work init interpreter compiler

go work sync
```

## summary

the language is called `Monkey` and has the following features：

- `C-like` syntax
- variable bindings
- integers and boolean
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures
- a string data structure
- an array data structure
- a hash data structure

```javascript
// here is how we bind values to names in Monkey
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2);

// array
let myArray = [1, 2, 3, 4, 5];

// hash
let thorsten = {"name": "Thorsten", "age": 28};

// accessing the elements in arrays and hashes is done with index expression
let intValue = myArray[0];
let name = thorsten["name"];

// the let statements can also be used to bind functions to names
let add = fn(a, b) {
  return a + b;
};

// implicit return values are also possible
let add_ = fn(a, b) {
  a + b
};

// a more complex function
let fibonacci = fn(x) {
  if (x == 0) {
    return 0;
  };
  if (x == 1) {
    return 1;
  };
  return fibonacci(x - 1) + fibonacci(x - 2);
};

fibonacci(10);

// a special type of functions, called higher order functions
let twice = fn(f, x) {
  return f(f(x));
};

let addTwo = fn(x) {
  return x + 2;
};

twice(addTwo, 2);
```

## module

the interpreter will have a few major parts:

- the lexer
- the parser
- the Abstract Syntax Tree(AST)
- the internal object system
- the evaluator
- the compiler
- the virtual machine

## structure

### Node

```mermaid
---
title: monkey AST
---
classDiagram
class Node
<<interface>> Node
class Expression
<<interface>> Expression
class Statement
<<interface>> Statement
Node <|-- Expression
Node <|-- Statement
Node <|-- Program
```

### Statement

```mermaid
---
title: monkey AST
---
classDiagram
class Statement
<<interface>> Statement

Statement <|-- BlockStatement
Statement <|-- ExpressionStatement
Statement <|-- LetStatement
Statement <|-- ReturnStatement
```

### Expression

```mermaid
---
title: monkey AST
---
classDiagram
class Expression
<<interface>> Expression

Node <|-- Expression

Expression <|--Array
Expression <|--Boolean
Expression <|--Call
Expression <|--FnLiteral
Expression <|--Hash
Expression <|--Identifier
Expression <|--IndexE
Expression <|--If
Expression <|--Infix
Expression <|--IntegerLiteral
Expression <|--Prefix
Expression <|--StringLiteral
```

## Recursive descent

According to structure, the `recursive descent` should be like the following code :

```go
func (c *Compiler) Compile(root ast.Node) error {
	switch node := root.(type) {
	case *ast.Program:
		return c.compileProgram(node)
	case ast.Expression:
		return c.compileExpression(node)
	// TODO support more ast
	case ast.Statement:
		return c.compileStatement(node)
	}
	return nil
}

func (c *Compiler) compileProgram(program *ast.Program) error {
	for _, stmt := range program.Statements {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return c.compileIntegerLiteral(expr)
		// TODO support more expression type
	case *ast.InfixExpression:
		return c.compileInfixOperator(expr)
	}
	return common.NewErrUnsupportedCompilingNode(expr.String())
}

func (c *Compiler) compileStatement(statement ast.Statement) error {
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		return c.compileExpressionStatement(stmt)
		// TODO support more statement type
	}

	return common.NewErrUnsupportedCompilingNode(statement.String())
}
```

## interpreter

```mermaid
mindmap
  root((monkey))
    token
    	Identifier
    	Operator
    	Delimiters
    	Preserved Keywords
    ast{{ast}}
    	Node
    		Program["**Program**
    			Program is the **root** of our program."]
    		Statement["**Statement**
    		combine **Node**, has a statementNode() method to specify a Node is a Statement"]
    			LetStatement
    			ReturnStatement
    		Expression["**Expression**
    		combine **Node**, has a expressionNode() method to specify a Node is a Expression"]
    			Identifier
    			ExpressionStatement
    lexer
    	NextToken["NextToken()
    	get next token from source code according current char and peek char"]
    	readChar["readChar()"]
    	peakChar["peakChar()"]
    parser
    	Top Down Operator Precedence
    	ParseProgram["ParseProgram()"]
    	BNF/EBNF
    	Pratt Parser
    		prefix operator
    		infix operator
    		suffix operator
    evaluator
    	Eval
    	Environment
    Object
    	Integer
    	Boolean
```

## compiler

### Jumps

```js
if (0 > 1) {
  2 + 3
} else {
  4 - 5
}
```

```mermaid
---
title: overview
---
block-beta
columns 4

0000:1 OpConstant0:3
0001:1 OpConstant1:3
0002:1 OpGreaterThan:3

0003:1 jump_if_not_true["JUMP_IF_NOT_TRUE"]:2 TO_0008["0008"]


0004:1 OpConstant2:3
0005 OpConstant3:3
0006 OpAdd:3

0007 jump_no_matter_what["JUMP_NO_MATTER_WHAT"]:2 TO_0011["0011"]

0008 OpConstant4:3
0009 OpConstant5:3
0010 OpMinus:3

0011 Code["..."]:3

jump_if_not_true --> 0008
jump_no_matter_what --> 0011

classDef front 1,fill:#696,stroke:#333;
classDef back fill:#969,stroke:#333;
classDef jump_to fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5

class 0000 front
class 0001 front
class 0002 front
class 0003 front
class 0004 front
class 0005 front
class 0006 front
class 0007 front
class 0008 front
class 0009 front
class 0010 front
class 0011 front

class jump_if_not_true back
class jump_no_matter_what back

class TO_0008 jump_to
class TO_0011 jump_to
```

### global identifiers

```js
let x = 33;
let y = 66;

let z = x + y;
```

```mermaid
---
title: global identifier
---
block-beta
columns 6

OFFSET:1 Instruction:1 Operand:2 Description:2
0000:1 Op0["OpConstant"]:1 Operand0["0"]:2 Desc0["Load the '33' onto the stack"]:2
0003:1 Op1["OpSetGlobal"]:1 Operand1["0"]:2 Desc1["Bind value on stack to 0"]:2
0006:1 Op2["OpConstant"]:1 Operand2["1"]:2 Desc2["Load the '66' onto the stack"]:2
0009:1 Op3["OpSetGlobal"]:1 Operand3["1"]:2 Desc3["Bind value on stack to 1"]:2

0012:1 Op4["OpGetGlobal"]:1 Operand4["1"]:2 Desc4["Push the global bound to 1"]:2
0015:1 Op5["OpGetGlobal"]:1 Operand5["0"]:2 Desc5["Push the global bound to 0"]:2
0018:1 Op6["OpAdd"]:3 esc6["Add them together"]:2
0019:1 Op7["OpSetGlobal"]:1 Operand7["2"]:2 Desc7["Bind value on stack to 1"]:2


classDef front 1,fill:#696,stroke:#333;
classDef back fill:#969,stroke:#333;
classDef op fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
classDef header fill: #696,color: #fff,font-weight: bold,padding: 10px;

class 0000 front
class 0003 front
class 0006 front
class 0009 front
class 0012 front
class 0015 front
class 0018 front
class 0019 front

class Op0 back
class Op1 back
class Op2 back
class Op3 back
class Op4 back
class Op5 back
class Op6 back
class Op7 back

class Operand0 op
class Operand1 op
class Operand2 op
class Operand3 op
class Operand4 op
class Operand5 op
class Operand6 op
class Operand7 op

class OFFSET header
class Instruction header
class Operand header
class Description header
```

The compiler is responsible for mapping identifiers to indices on the stack. And in the VM we'll use a slice to implement the creation and retrieval of global bindings.We'll call this slice our "global store" and we'll use the operands of the `OpSetGlobal` and `OpGetGlobal` instructions as indexes into it.

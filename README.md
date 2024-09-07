# monkey

> **writing an interpreter/compiler in Go**
>
> In this repository, I am going to build my own lexer, my own parser, my own tree representation and my own evaluator.

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
    		combine **Node**, has a statementNode() method to speicify a Node is a Statement"]
    			LetStatement
    			ReturnStatement
    		Expression["**Expression**
    		combine **Node**, has a expressionNode() method to speicify a Node is a Expression"]
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
    		sufix operator
    evaluater
    	Eval
    	Environment
    Object
    	Integer
    	Boolean
```


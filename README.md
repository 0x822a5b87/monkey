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

## usage

### primary type, hash and array

```js
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
```

### function

```js
// the let statements can also be used to bind functions to names
let add = fn(a, b) {
  return a + b;
};

// implicit return values are also possible
let add_ = fn(a, b) {
  a + b
};
```

### recursive

```js
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
```

### closure

#### a simple closure with first-class function

```js
let returnsOne = fn() { 1; };
let returnsOneReturner = fn() { returnsOne; };

let closure = returnsOneReturner();
closure()
```

#### a more complex closure

```javascript
// a special type of functions, called higher order functions
let twice = fn(f, x) {
  return f(f(x));
};

let addTwo = fn(x) {
  return x + 2;
};

twice(addTwo, 2);
```

### local bindings

```js
let globalSeed = 50;

let minusOne = fn() {
  let num = 1;
  globalSeed - num;
};

let minusTwo = fn() {
  let num = 2;
  globalSeed - num;
};

// 97
minusOne() + minusTwo();
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

### functions

consider the follwing code snippet: A function without arguments, an integer arithmetic expression in the body, and a *explicit return statement*.

```js
fn() {
  return 5 + 10;
}
```

```mermaid
---
title: functions with explicit return statement
---
block-beta
columns 6

OFFSET:1 Instruction:1 Operand:2 Description:2
0000:1 Op0["OpConstant"]:1 Operand0["0"]:2 Desc0["Load 5 on stack"]:2
0003:1 Op1["OpConstant"]:1 Operand1["1"]:2 Desc1["Load 10 on stack"]:2
0006:1 Op2["OpAdd"]:3 Desc2["Add them together"]:2
0007:1 Op3["OpReturnValue"]:3 Desc3["return value on stack"]:2

classDef front 1,fill:#696,stroke:#333;
classDef back fill:#969,stroke:#333;
classDef op fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
classDef header fill: #696,color: #fff,font-weight: bold,padding: 10px;

class 0000 front
class 0003 front
class 0006 front
class 0007 front

class Op0 back
class Op1 back
class Op2 back
class Op3 back

class Operand0 op
class Operand1 op
class Operand2 op
class Operand3 op

class OFFSET header
class Instruction header
class Operand header
class Description header
```

### frames

Consider this:

```js
// define
let one = fn() { 5; };
let two = fn() { one(); };
let three = fn() { two(); };

// call
three();
```

As we know, function calls are nested and execution-relevant data(the instructions and the instruction pointer) is accessed in a last-in-first-out manner. The solution is to tie them together and call the resulting bundle a `frame` -- short for `call frame` or `stack frame`.

On real mechines, a frame is not something separate from but a designated part of the stack.**It's where the return address, the arguments to the current function and its local variables are stored.**

```go
package vm

import "0x822a5b87/monkey/compiler/code"

// Frame short for call frame or stack frame
// A Frame has two fields: ip and fn.
// fn points to the compiled function referenced by the frame.
// ip is the instruction pointer in this frame, for the function.
type Frame struct {
	fn *code.CompiledFunction
	ip int
}

func NewFrame(f *code.CompiledFunction) *Frame {
	return &Frame{
		fn: f,
		ip: -1,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}

```

### local bindings on stack

When we come across an `OpCall` instruction in the VM and are about to execute the function on stack, we take the current value of the stack pointer and put it aside -- for later use.We then increase the stack pointer by the number of locals used by the function we're about to execute. The result is a "hole" on the stack: we've increased the stack pointer without pushing any values, creating a region on the stack without any values.Below the hole: all the values previously pushed on the stack, before the function call. And above the hole is the functions workspace, where it will push and pop the values it needs to do its work.

The hole itself is where we're going to store local bindings. We won't use the unique index of a local bindings as a key for another data structure, but instead as an index into the hole on the stack.

>The entire execution process of the function is as follows:
>
>1. Encountering the `OpConstant` instruction or `OpGetGlobal` instruction, which retrieves the CompiledFunction from the stack.
>2. Encountering the OpCall instruction, which executes the result obtained from the execution of OpConstant.
>3. Saving the stack pointer of the current frame into the new frame's context, which will be used to reset the stack pointer after the function call ends.
>4. Entering the new frame;
>5. Allocating the necessary positions for local variables and function;
>6. Executing the function;
>7. Resetting stack pointer to the saved value;
>8. Exiting the current frame;

```mermaid
---
title: function stack
---
block-beta
columns 3

block:group1:2
	columns 1
	stack3["stack 3"]:1
	stack2["stack 2"]:1
	stack1["stack 1"]:1
	stack0["stack 0"]:1
end
desc1["empty"]
group1 --> desc1

block:group2:2
	columns 1
	Local3["..."]:1
	Local2["Local 2"]:1
	Local1["Local 1"]:1
	Function
end
desc2["reserved for locals"]
group2 --> desc2

block:group3:2
	columns 1
	OtherValue1["..."]:1
	OtherValue2["Other Value 3"]:1
	OtherValue3["Other Value 2"]:1
	OtherValue4["Other Value 1"]:1
end
desc3["pushed before call"]
group3 --> desc3

classDef front 1,fill:#696,stroke:#333;
classDef back fill:#969,stroke:#333;
classDef op fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
classDef header fill: #696,color: #fff,font-weight: bold,padding: 10px;

class stack3 op
class stack2 op
class stack1 op
class stack0 op

class Local3 header
class Local2 header
class Local1 header
class Function header

class OtherValue1 back
class OtherValue2 back
class OtherValue3 back
class OtherValue4 back
```

### Arguments

> arguments to function calls are a special case of local bindings.
>
> The have the same lifespan, the have the same scope, they resolve in the same way.The only different is their creation. 
>
> - Local bindings are created explicitly by the user with a let statement and result in `OpSetLocal` instructions being emitted by the compiler.
> - Arguments are implicitly bound to names, which is done behind the scenes by the compiler and the VM.

There are three essential pointer in the implementation called `stackPointer`, `basePointer`, `functionPointer`:

1. base pointer points to the start position of local variables;
2. stack pointer points to the start position of the new frame's stack;
3. function pointer points to the position of function;

When encountering a new function call, we create a new frame with these three pointer:

1. Retrieve function through function pointer;
2. Execute statements of the function and place the result onto stack through stack pointer;
3. Get and set local variable through base pointer and relative offset of each variable, for example : `Local 0` == `basePointer + 0`. 
4. **When the execution of function is completed, we reset the stack pointer of the parent frame to function pointer which is the original value of the parent frame**

```mermaid
---
title: function stack
---
block-beta
columns 3

block:group1:2
	columns 2
	stack3["stack 3"]:2
	stack2["stack 2"]:2
	stack1["stack 1"]:2
	stack0["stack 0"]:1
    stackPointer:1
end
stackPointer --> stack0
desc1["empty"]:1

block:group2:2
	columns 2
	Local2["..."]:2
	Local1["Local 1"]:2
	Local0["Local 0"]:1
    basePointer
	Function:1
	functionPointer
end
desc2["reserved for locals"]:1
basePointer --> Local0
functionPointer --> Function

block:group3:2
	columns 1
	OtherValue1["..."]:1
	OtherValue2["Other Value 3"]:1
	OtherValue3["Other Value 2"]:1
	OtherValue4["Other Value 1"]:1
end
desc3["pushed before call"]

classDef front 1,fill:#FFCCCC,stroke:#333;
classDef back fill:#969,stroke:#333;
classDef op fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
classDef header fill: #696,color: #fff,font-weight: bold,padding: 10px;

class stack3 op
class stack2 op
class stack1 op
class stack0 op

class Local0 header
class Local2 header
class Local1 header
class Function header

class OtherValue1 back
class OtherValue2 back
class OtherValue3 back
class OtherValue4 back

class basePointer front
class stackPointer front
class functionPointer front
```

### Closures

Closures in summary:

1. detect references to free variables while compiling a function
2. get the referenced values on to the stack
3. **merge the values and the compiled function into a closure and leave it on the stack where it can then be called.**

More specifically, closure invovles two relative instructions:

| instruction | operands | description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| OpClosure   | 2        | 1. The first operand is the constant index, it specifies where in the constant pool we can find the CompiledFunction that's to be converted into a closure.\n2. The second operand, it specifies how many free variables sit on the stackand to be translated to about-to-be-created closure. |
| OpGetFree   | 1        | index of free variable in the the closure field              |
























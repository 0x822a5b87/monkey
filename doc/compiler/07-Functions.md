# Functions

> how to represent functions?

The naive assumptions is that functions are a series of instructions. But in Monkey functions are also first-class citizens that can passed around and returned from other functions. **How do we represent a series of instructions that can be passed around?**

> the issue of control flow

1. how do we get our VM to execute the instructions of a function?
2. how do we get it to **return back** to the instructions it was previously executing?
3. how do we get it to pass **arugments** to functions?

## Dipping Our Toes: a Simple Function

### Representing Functions

> We have found a ingenious way to handle the problems we encountered by compiling the function literal into an object.

```go
package code

import (
	"0x822a5b87/monkey/interpreter/object"
	"fmt"
)

const (
	ObjCompiledFunction object.ObjType = "COMPILED_FUNCTION"
)

// CompiledFunction a function object that holds bytecode instead of AST nodes.
// It can hold the Instructions we get from the compilation of a function literal, and it's an object.Object, which means
// we can add it as a constant to our compiler.ByteCode and load it in the VM
type CompiledFunction struct {
	Instructions Instructions
}

func (c *CompiledFunction) Type() object.ObjType {
	return ObjCompiledFunction
}

func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction [%p]", c.Instructions)
}
```

### Opcodes to Execute Functions

The first question we have to ask ourselves is whether we need new opcodes to achieve our goal of compiling and executing the snippet of Monkey code from above.Consider the following code:

```js
let fivePlusTen = fn() { 5 + 10 };

fivePlusTen();
```

1. We don't need an opcode for function literals.
2. We need an opcode for function call.

Once we have compiled the function literal into an `*object.CompiledFunction` we already know how to bind it to the *fivePlusTen* name.

And a call expression must be compiled into an instruction that tells the VM to execute the function in question.

With OpCall defined, we are now able to get a function on to the stack of our VM and call it. What we still need is a way to tell the VM to return from a called function. More specifically, we need to differentiate between two cases where the Vm has to return from a function:

1. the execution of a function actually returning something;
2. the execution of a function ends without anything being returned;
3. a much rarer case when returnning from a function: a function returning nothing, neither explicitly nor implicitly.

```js
// in this case, we will return a null value. We'll do that by introducing another opcode called OpReturn
fn() { }
```

### Compling Function Literals

Before we start opening our test file, a little inventory check. We now have in place:

- `object.CompiedFunction` to hold the instructions of a compied function and to pass them from the compiler to the VM as part of the bytecode, **as a constant**.
- `code.OpCall` to tell the VM to start executing the `*object.CompilerFunction` sitting on top of the stack.
- `code.OpReturnValue` to tell the VM to return the value on top of the stack to the calling context and to resume execution there.
- `code.OpReturn`, which is similar to `code.OpReturnValue`, except that there is no explicit value to return but an implici `Null`


























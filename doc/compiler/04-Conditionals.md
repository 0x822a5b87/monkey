# Conditionals

Let's say we have the following monkey code:

```javascript
if (5 > 2) {
  30 + 20
} else {
  50 - 25
}
```

> represent the condition `5 > 2` 

```mermaid
---
title: 5 > 2
---
flowchart TD

subgraph instructions
	direction TB
	lhs[OpConstant 0]
	rhs[OpConstant 1]
	operator[OpGreaterThan]
end

lhs --> stack0
rhs --> stack1

subgraph stack
	direction TB
	stack0[5]
	stack1[2]
	other[...]
end
```

> represent the consequence `30 + 20`

```mermaid
---
title: 30 + 20
---
flowchart TD

subgraph instructions
	direction TB
	lhs[OpConstant 2]
	rhs[OpConstant 3]
	operator[OpAdd]
end

lhs --> stack2
rhs --> stack3

subgraph stack
	direction TB
	stack0[5]
	stack1[2]
	stack2[30]
	stack3[20]
	other[...]
end
```

> represent the alternative `50 -25`

```mermaid
---
title: 50 - 25
---
flowchart TD

subgraph instructions
	direction TB
	lhs[OpConstant 4]
	rhs[OpConstant 5]
	operator[OpSub]
end

lhs --> stack4
rhs --> stack5

subgraph stack
	direction TB
	stack0[5]
	stack1[2]
	stack2[30]
	stack3[20]
	stack4[50]
	stack5[25]
	other[...]
end
```

> represent the whole code snippet

```mermaid
---
title: overview
---
flowchart TD

subgraph condition
	direction TB
	op0[OpConstant 0]
	op1[OpConstant 1]
	operator1[OpGreaterThan]
end

JUMP_IF_NOT_TRUE

subgraph consequence
	direction TB
	op2[OpConstant 2]
	op3[OpConstant 3]
	operator2[OpAdd]
end

subgraph alternative
	direction TB
	op4[OpConstant 4]
	op5[OpConstant 5]
	operator3[OpSub]
end

op0 --> stack0
op1 --> stack1
op2 --> stack2
op3 --> stack3
op4 --> stack4
op5 --> stack5

subgraph stack
	direction TB
	stack0[5]
	stack1[2]
	stack2[30]
	stack3[20]
	stack4[50]
	stack5[25]
	other[...]
end
```

If we were to take these instructions and pass them to the VM as a flat sequence, what would happen?The VM would execute all of them, one after the other, happily incrementing its instruction pointer, fetching, decoding and executing.And that's exactly what we don't want!

we need put something in the blanks so that based on the result of the `OpGreaterThan` instrction the VM either ignores the instrcutions of the consequence or the instructions making up the alternative. It should skip them. Or instead of `skip`, should we maybe say `jump over`?

### Jumps

`Jumps` are instructions that tell machines to jump to other instructions, or we can say: **jumps are instructions that tell the VM to change its instruction pointer to a certain value.**

```mermaid
---
title: overview
---
block-beta
columns 4

0000:1 OpConstant0:3
0001:1 OpConstant1:3
0002:1 OpGreaterThan:3

0003:1 jump_if_not_true["JUMP_IF_NOT_TRUE"]:3

0004:1 OpConstant2:3
0005 OpConstant3:3
0006 OpAdd:3

0007 jump_no_matter_what["JUMP_NO_MATTER_WHAT"]:3

0008 OpConstant4:3
0009 OpConstant5:3
0010 OpMinus:3

0011 Code["..."]:3

jump_if_not_true --> 0008
jump_no_matter_what --> 0011

classDef front 1,fill:#696,stroke:#333;
classDef back fill:#969,stroke:#333;
  
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
```

we use numbers to tell the VM where to jump to, and it can be seperated to different type: `absolute offset` and `relative offset`.

- absolute offset : the jump target being the index of an jump instruction
- relative offset : relative to the position of the jump instruction itself and denote not where exactly to jump to, but how far to jump.

If we represent the code with offsets and give each instruction a unique index that's indpendent of its byte size(for inllustration purposes), the diagram looks like this:

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










































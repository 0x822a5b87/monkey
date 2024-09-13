# Keep track of names

```js
let x = 10;
```

## compile

compileExpression

1. push integer constant value 10 to constant pool;
2. generate instruction : OpConstant constantIndex

compileLetStatement

1. define variable in symbol tables
2. generate instruction : OpSetGlobal globalIndex

## run

OpConstant constantIndex

1. read operand which is index of integer constant value 10
2. push index to the topmost stack

OpSetGlobal globalIndex

1. pop topmost value off the stack
2. set global value using globalIndex as key and the value comes from pop as value

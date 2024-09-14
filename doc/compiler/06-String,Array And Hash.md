# String, Array and Hash

In the end, we can exuecute this piece of monkey code :

```js
[1, 2, 3][1]
// => 2

{"one": 1, "two": 2, "three": 3}["o" + "ne"]
// => 1
```

## String

Since the value of string literals doesn't change between compile and runtime, we can treat them as constant expression.Similar to our implementation of integer literals.

Instead of only implementing strings literals, we'll also make it a goal for this section to implement string concatenation, which allows us to concatenate two strings with the `+` operator.


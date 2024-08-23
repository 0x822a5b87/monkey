package object

// Subtract the behavior of the infix operator minus
type Subtract interface {
	Object
	Sub(Object) Object
}

// Multiply the behavior of the infix operator asterisk
type Multiply interface {
	Object
	Mul(Object) Object
}

type Equatable interface {
	Object
	Equal(Object) *Boolean
	NotEqual(Object) *Boolean
}

type Comparable interface {
	Object
	GreaterThan(Object) *Boolean
	LessThan(Object) *Boolean
}

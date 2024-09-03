package object

// Add the operation of the infix operator +
type Add interface {
	Object
	Add(Object) Object
}

// Subtract the operation of the infix operator -
type Subtract interface {
	Object
	Sub(Object) Object
}

// Multiply the operation of the infix operator *
type Multiply interface {
	Object
	Mul(Object) Object
}

// Divide the operation of the infix operator /
type Divide interface {
	Object
	Divide(object Object) Object
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

type Negative interface {
	Object
	Negative() Object
}

type Index interface {
	Object
	Index(Object) Object
	First() Object
	Last() Object
}

type Len interface {
	Object
	Len() Integer
}

type Rest interface {
	Object
	Rest() Object
}

type Push interface {
	Object
	Push(obj Object) Object
}

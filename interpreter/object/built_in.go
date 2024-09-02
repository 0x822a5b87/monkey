package object

type BuiltInFunction func(objs ...Object) Object

type BuiltIn struct {
	BuiltInFn BuiltInFunction
}

func (b *BuiltIn) Type() ObjType {
	return ObjBuiltIn
}

func (b *BuiltIn) Inspect() string {
	return "built-in function"
}

type BuiltInFn struct {
}

package object

const (
	builtInFnNameLet   = "len"
	builtInFnNameFirst = "first"
	builtInFnNameLast  = "last"
	builtInFnNameRest  = "rest"
	builtInFnNamePush  = "push"
)

var id int

var globalEnv *Environment

func NewEnvironment(parent *Environment) *Environment {
	id++
	if parent == nil {
		parent = globalEnv
	}
	return &Environment{
		name:   id,
		store:  make(map[string]Object),
		parent: parent,
	}
}

type Environment struct {
	name   int
	store  map[string]Object
	parent *Environment
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.parent != nil {
		obj, ok = env.parent.Get(name)
	}
	return obj, ok
}

func (env *Environment) Set(name string, obj Object) {
	env.store[name] = obj
}

func init() {
	globalEnv = &Environment{
		name:   0,
		store:  make(map[string]Object),
		parent: nil,
	}

	globalEnv.Set(builtInFnNameLet, &BuiltIn{BuiltInFn: func(objs ...Object) Object {
		if len(objs) != 1 {
			return newWrongArgumentSizeError(len(objs), 1)
		}

		obj := objs[0]
		builtInLen, ok := obj.(Len)
		if !ok {
			return newWrongArgumentTypeError(builtInFnNameLet, obj.Type())
		}
		return &Integer{Value: builtInLen.Len().Value}
	}})

	globalEnv.Set(builtInFnNameFirst, &BuiltIn{BuiltInFn: func(objs ...Object) Object {
		if len(objs) != 1 {
			return newWrongArgumentSizeError(len(objs), 1)
		}

		obj := objs[0]
		builtInFirst, ok := obj.(Index)
		if !ok {
			return newWrongArgumentTypeError(builtInFnNameFirst, obj.Type())
		}
		return builtInFirst.First()
	}})

	globalEnv.Set(builtInFnNameLast, &BuiltIn{BuiltInFn: func(objs ...Object) Object {
		if len(objs) != 1 {
			return newWrongArgumentSizeError(len(objs), 1)
		}

		obj := objs[0]
		buildInLast, ok := obj.(Index)
		if !ok {
			return newWrongArgumentTypeError(builtInFnNameLast, obj.Type())
		}
		return buildInLast.Last()
	}})

	globalEnv.Set(builtInFnNameRest, &BuiltIn{BuiltInFn: func(objs ...Object) Object {
		if len(objs) != 1 {
			return newWrongArgumentSizeError(len(objs), 1)
		}

		obj := objs[0]
		rest, ok := obj.(Rest)
		if !ok {
			return newWrongArgumentTypeError(builtInFnNameRest, obj.Type())
		}
		return rest.Rest()
	}})

	globalEnv.Set(builtInFnNamePush, &BuiltIn{BuiltInFn: func(objs ...Object) Object {
		if len(objs) != 2 {
			return newWrongArgumentSizeError(len(objs), 2)
		}

		obj := objs[0]
		push, ok := obj.(Push)
		if !ok {
			return newWrongArgumentTypeError(builtInFnNamePush, obj.Type())
		}
		return push.Push(objs[1])
	}})
}

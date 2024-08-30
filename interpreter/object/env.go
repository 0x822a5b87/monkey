package object

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		store:  make(map[string]Object),
		parent: parent,
	}
}

type Environment struct {
	store  map[string]Object
	parent *Environment
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	for !ok && env.parent != nil {
		obj, ok = env.parent.Get(name)
	}
	return obj, ok
}

func (env *Environment) Set(name string, obj Object) {
	env.store[name] = obj
}

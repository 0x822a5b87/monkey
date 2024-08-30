package object

var id int

func NewEnvironment(parent *Environment) *Environment {
	id++
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

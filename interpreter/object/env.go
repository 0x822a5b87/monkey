package object

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

type Environment struct {
	store map[string]Object
}

func (env *Environment) Get(name string) (Object, bool) {
	o, ok := env.store[name]
	return o, ok
}

func (env *Environment) Set(name string, obj Object) {
	env.store[name] = obj
}

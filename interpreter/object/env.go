package object

const (
	builtInFnNameLen   = "len"
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

	for _, builtIn := range BuiltIns {
		globalEnv.Set(builtIn.Name, builtIn)
	}
}

package compiler

// A symbol table is a data structure used in interpreters and compilers to associate identifiers
// with information. It can be used in every phase, from lexing to code generation,
// to store and retrieve information about a given identifiers(which can be called a symbol).

type SymbolScope string

const (
	GlobalScope SymbolScope = "Global"
	LocalScope  SymbolScope = "Local"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable associates strings with Symbol in its store and keeps track of the numDefinitions it has.
type SymbolTable struct {
	Outer *SymbolTable
	store map[string]Symbol
	// TODO figure out its usage
	// num of this symbol table have been defined which used to prevent it being release by GC at an inappropriate time
	numDefinitions int
}

func (st *SymbolTable) Define(name string) Symbol {
	scope := GlobalScope
	if st.Outer != nil {
		scope = LocalScope
	}

	st.checkDefine()
	s := Symbol{
		Name:  name,
		Scope: scope,
		Index: st.numDefinitions,
	}
	st.store[name] = s
	st.numDefinitions++
	return s
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	if !ok && st.Outer != nil {
		return st.Outer.Resolve(name)
	}
	return s, ok
}

func (st *SymbolTable) checkDefine() bool {
	// TODO check if the definition is legal
	return true
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func NewEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	enclosed := NewSymbolTable()
	enclosed.Outer = parent
	return enclosed
}

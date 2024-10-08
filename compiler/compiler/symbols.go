package compiler

import "0x822a5b87/monkey/interpreter/object"

// A symbol table is a data structure used in interpreters and compilers to associate identifiers
// with information. It can be used in every phase, from lexing to code generation,
// to store and retrieve information about a given identifiers(which can be called a symbol).

type SymbolScope string

const (
	GlobalScope  SymbolScope = "Global"
	LocalScope   SymbolScope = "Local"
	BuiltInScope SymbolScope = "BuiltIn"
	FreeScope    SymbolScope = "Free"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable associates strings with Symbol in its store and keeps track of the numDefinitions it has.
type SymbolTable struct {
	Outer *SymbolTable
	// Free contains the original value of free variables
	// it's important to know free variables also have its own scope, and it exits in a symbol table like regular variables
	Free  []Symbol
	store map[string]Symbol
	// num of this symbol table have been defined which used to prevent it being release by GC at an inappropriate time
	numDefinitions int
	numFree        int
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

func (st *SymbolTable) DefineBuiltIn(index int, name string) Symbol {
	st.checkDefine()
	s := Symbol{
		Name:  name,
		Scope: BuiltInScope,
		Index: index,
	}
	st.store[name] = s
	return s
}

// Resolve a symbol with a local scope, global scope, built-in function, or free scope through its name.
// For local or global variables, we only need to check if the store contains the corresponding symbol.
// For free variables, we have to perform the following checks:
// 1. Has the name been defined in **this** scope, **this** symbol table ? If true, return a symbol with regular scope.
// 2. Is it a global binding or a built-in function? If true, return a symbol with regular scope.
// If none of the previous conditions are met, it indicates the variable is a free variable.
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	// check this scope
	s, ok := st.store[name]

	// check built-in function
	if !ok && st.Outer != nil {
		s, ok = st.Outer.Resolve(name)
		// failed to resolve symbol
		if !ok {
			return s, ok
		}

		// test if the variable is free variable or built-in function
		if s.Scope != BuiltInScope && s.Scope != GlobalScope {
			s = st.defineFree(s)
		}
	}

	return s, ok
}

func (st *SymbolTable) checkDefine() bool {
	// TODO check if the definition is legal
	return true
}

func (st *SymbolTable) defineFree(original Symbol) Symbol {
	st.checkDefine()
	st.Free = append(st.Free, original)
	symbol := Symbol{Name: original.Name, Index: len(st.Free) - 1}
	symbol.Scope = FreeScope
	st.store[original.Name] = symbol
	st.numFree++
	return symbol
}

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{
		Free:  make([]Symbol, 0),
		store: make(map[string]Symbol),
	}
	return s
}

func NewGlobalSymbolTable() *SymbolTable {
	s := NewSymbolTable()
	for index, builtIn := range object.BuiltIns {
		s.DefineBuiltIn(index, builtIn.Name)
	}
	return s
}

func NewEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	enclosed := NewSymbolTable()
	enclosed.Outer = parent
	return enclosed
}

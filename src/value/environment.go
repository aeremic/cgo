package value

type Environment struct {
	store map[string]Wrapper
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Wrapper)

	return &Environment{
		store: s,
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func (e *Environment) Get(name string) (Wrapper, bool) {
	wrappedValue, ok := e.store[name]
	if !ok && e.outer != nil {
		wrappedValue, ok = e.outer.Get(name)
	}

	return wrappedValue, ok
}

func (e *Environment) Set(name string, wrappedValue Wrapper) Wrapper {
	e.store[name] = wrappedValue
	return wrappedValue
}

package value

type Environment struct {
	store map[string]Wrapper
}

func NewEnvironment() *Environment {
	s := make(map[string]Wrapper)

	return &Environment{
		store: s,
	}
}

func (e *Environment) Get(name string) (Wrapper, bool) {
	wrappedValue, ok := e.store[name]

	return wrappedValue, ok
}

func (e *Environment) Set(name string, wrappedValue Wrapper) Wrapper {
	e.store[name] = wrappedValue
	return wrappedValue
}

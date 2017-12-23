package fit

type Shared struct {
	values map[string]interface{}
}

func (c *Context) Shared() Shared {
	if c.shared.values == nil {
		c.shared = Shared{make(map[string]interface{})}
	}

	return c.shared
}

func (s Shared) Get(key string) (ok bool, value interface{}) {
	value, ok = s.values[key]
	return
}

func (s Shared) Set(key string, value interface{}) {
	s.values[key] = value
}

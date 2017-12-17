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

func (s Shared) Get(key string) (bool, interface{}) {
	if k, ok := s.values[key]; ok {
		return true, k
	}

	return false, ""
}

func (s Shared) Set(key string, value interface{}) {
	s.values[key] = value
}

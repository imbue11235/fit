package fit

import "fmt"

type Shared struct {
	context *Context
}

func (c *Context) Shared() Shared {
	if c.shared == nil {
		c.shared = make(map[string]interface{})
	}

	return Shared{c}
}

func (s Shared) Get(key string) (bool, interface{}) {
	if k, ok := s.context.shared[key]; ok {
		return true, k
	}

	fmt.Printf("Key %s does not exist. Ignoring.\n", key)

	return false, ""
}

func (s Shared) Set(key string, value interface{}) {
	s.context.shared[key] = value
}

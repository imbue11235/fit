package fit

type Parameters struct {
	stack []parameter
}

type parameter struct {
	key   string
	value string
}

// Parameters ....
func (c *Context) Parameters() Parameters {
	return c.params
}

// GetByName ....
func (p Parameters) GetByName(name string) (bool, string) {
	for _, parameter := range p.stack {
		if parameter.key == name {
			return true, parameter.value
		}
	}
	return false, ""
}

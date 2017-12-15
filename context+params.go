package fit

type Params struct {
	params map[string]string
}

// Parameters ....
func (c *Context) Parameters() Params {
	return c.params
}

// GetByName ....
func (p Params) GetByName(name string) (bool, string) {
	if param, ok := p.params[name]; ok {
		return true, param
	}
	return false, ""
}

// Get by int?

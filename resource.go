package fit

type resource struct {
	path     string
	methods  map[string][]ResponseHandler
	prefix   string
	children []*resource
	options  *Options
}

// Helper functions
func newResource() *resource {
	return &resource{
		path:     "",
		methods:  make(map[string][]ResponseHandler),
		children: make([]*resource, 0),
		prefix:   "",
		options:  &Options{},
	}
}

func newResourceFromPath(path string) *resource {
	res := newResource()
	res.path = path
	return res
}

func (res *resource) copy() *resource {
	cop := new(resource)
	*cop = *res
	return cop
}

func (res *resource) addMethods(methods []string, options *Options, handlers ...ResponseHandler) {
	for _, m := range methods {
		if _, ok := res.methods[m]; ok {
			panic("handler existed!")
		}
		res.methods[m] = handlers
		res.options = options
	}
}

func (res *resource) getIndexPosition(target byte) int {
	min, max := 0, len(res.prefix)
	for min < max {
		mid := min + ((max - min) >> 1)
		if res.prefix[mid] < target {
			min = mid + 1
		} else {
			max = mid
		}
	}
	return min
}

func (res *resource) insertChild(index byte, child *resource) *resource {

	i := res.getIndexPosition(index)

	res.prefix = res.prefix[:i] + string(index) + res.prefix[i:]
	res.children = append(res.children[:i], append([]*resource{child}, res.children[i:]...)...)

	return child
}

func (res *resource) getChild(index byte) *resource {
	i := res.getIndexPosition(index)
	if i == len(res.prefix) || res.prefix[i] != index {
		return nil
	}
	return res.children[i]
}

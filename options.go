package fit

import (
	"fmt"
	"strings"
)

type Options struct {
	regex map[string]string
	name  string
	path  string
}

// Name - Sets name of the route
func (r *Options) Name(name string) *Options {
	r.name = name
	return r
}

// Middleware ...
/*
func (r *Options) Middleware(middleware ...ContextFunc) {
	r.middleware = append(r.middleware, middleware...)
}
*/

// Where ...
func (r *Options) Where(constraints ...string) *Options {
	regex, constraintLength := r.regex, len(constraints)
	if regex == nil {
		regex = make(map[string]string)
	}

	if constraintLength%2 != 0 {
		// TODO: What does this mean?
		panic("Constraint is missing")
	}

	for i := 0; i < constraintLength; i += 2 {
		constraintName, constraintValue := constraints[i], constraints[i+1]

		// Empty constraint(s)
		if constraintValue == "" || constraintName == "" {
			fmt.Println("Empty constraint was supplied. Ignoring")
			continue
		}

		// If constraint does not exist in url path, we will ignore it
		if !strings.Contains(r.path, constraintName) {
			fmt.Printf("Constraint '%s' does not exist in path '%s'. Ignoring.\n", constraintName, r.path)
			continue
		}
		regex[constraintName] = constraintValue
	}

	r.regex = regex
	return r
}

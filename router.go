package fit

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

const (
	star  = byte('*')
	colon = byte(':')
	slash = byte('/')
)

type Router struct {
	res             *resource
	logger          ResponseHandler
	IgnoreFavicon   bool
	RedirectSlashes bool
}

func NewRouter() *Router {
	return &Router{newResource(), nil, true, true}
}

// Serve ....
func (r *Router) Serve(port ...int) {
	// Default port
	portString := ":8080"
	if len(port) > 0 {
		portString = fmt.Sprintf(":%d", port[0])
	}

	if r.IgnoreFavicon {
		// Do ignore requests to the favicon. Not really needed in Rest API's.
		// I've added the option to choose to ignore it or not. Default is ignoring.
		http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {}) // Ignore favicon
	}

	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		found, handlers, params := r.findRoute(rq.URL.Path, rq.Method)

		c := newContext()
		c.writer = w
		c.request = rq

		if found && len(handlers) > 0 {
			// CREATE THIS SOMEWHERE ELSE

			c.params = Params{params}
			c.handlers = handlers
			c.currentHandler = 0
			c.maxHandlers = len(handlers)

			c.callByIndex(0)
		} else {
			c.status = http.StatusNotFound
			// Error handler here
			fmt.Fprintln(w, "Does not exist. Should return appropriate status 404")

		}
		r.logger(c)

	})

	log.Fatal(http.ListenAndServe(portString, nil))
}

func (r *Router) Logger(handler ResponseHandler) {
	r.logger = handler
}

func (r *Router) addRoute(path string, methods []string, handlers ...ResponseHandler) *Options {
	i, pathLength, res, options := 0, len(path), r.res, &Options{path: path}

	for i < pathLength {
		position := res.getIndexPosition(path[i])
		if position == len(res.prefix) || res.prefix[position] != path[i] {
			position = find(path, colon, i, pathLength)
			if position == pathLength {
				position = find(path, star, i, pathLength)
				res = res.insertChild(path[i], newResourceFromPath(path[i:position]))

				if position < pathLength {
					res = res.insertChild(star, newResourceFromPath(path[position+1:]))
				}

				res.addMethods(methods, options, handlers...)
				break
			}

			res = res.insertChild(path[i], newResourceFromPath(path[i:position]))
			i = find(path, slash, position, pathLength)
			res = res.insertChild(colon, newResourceFromPath(path[position+1:i]))

			if i == pathLength {
				res.addMethods(methods, options, handlers...)
			}

		} else if path[i] == colon {
			res = res.children[0]
			i += len(res.path) + 1

			if i == pathLength {
				res.addMethods(methods, options, handlers...)
			}

		} else {
			res = res.getChild(path[i])

			j, resourcePathLength := 0, len(res.path)
			for j < resourcePathLength && i < pathLength && path[i] == res.path[j] {
				i++
				j++
			}

			if j < resourcePathLength {

				child := res.copy()
				child.path = res.path[j:]

				// Why cant i simplify this with newResource?
				res.path = res.path[:j]
				res.methods = make(map[string][]ResponseHandler)
				res.prefix = string(child.path[0])
				res.children = []*resource{child}
				res.options = &Options{}
			}

			if i == pathLength {
				res.addMethods(methods, options, handlers...)
			}
		}
	}
	return options
}

func (r *Router) findRoute(path, method string) (found bool, handlers []ResponseHandler, params map[string]string) {

	i, pathLength, res, params := 0, len(path), r.res, make(map[string]string)

	for i < pathLength {
		if len(res.prefix) == 0 {
			return
		}

		if res.prefix[0] == colon {
			res = res.children[0]
			position := find(path, slash, i, len(path))
			params[res.path] = path[i:position]
			i = position
		} else if res.prefix[0] == star {
			res = res.children[0]
			params[res.path] = path[i:]
			break
		} else {
			position := res.getIndexPosition(path[i])
			if position == len(res.prefix) || res.prefix[position] != path[i] {
				return
			}
			res = res.children[position]
			position = i + len(res.path)
			if position >= pathLength || path[i:position] != res.path {
				return
			}
			i = position
		}
	}

	// If regex is specified, we will run it against the parameters
	if res.options.regex != nil {
		for name, constraint := range res.options.regex {
			if param, ok := params[name]; ok {
				validRoute := regexp.MustCompile(constraint)
				if !validRoute.MatchString(param) {
					// Not found
					return
				}
			}
		}
	}

	return true, res.methods[method], params
}

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

// Router ...
type Router struct {
	// Resource tree for the assigned routes
	res *resource

	// Contains ResponseHandler(s) called before the handlers assigned on the route
	before []ResponseHandler

	// Contains ResponseHandler(s) called after the handlers assigned on the route
	after []ResponseHandler

	// Contains a ResponseHandler called after everything, not dependent of the middleware chain
	logger ResponseHandler

	// Contains the default function to use when a page was not found (404)
	NotFound ResponseHandler

	// A boolean value for toggling automatic redirects, if the route exists with (or without) slashes "/"
	RedirectSlashes bool
}

// NewRouter returns a new instance of the Router struct.
// It's created with an empty resource and a standard not found handler for 404 requests.
func NewRouter() *Router {
	return &Router{
		newResource(),     // Resource creation
		nil,               // Before ResponseHandler(s)
		nil,               // After ResponseHandler(s)
		nil,               // Logger ResponseHandler
		notFoundHandler(), // Default not found handler
		true,              // RedirectSlashes is activated pr. default
	}
}

// Before appends handler(s) before all other handlers, globally for the instance of the router
func (r *Router) Before(handlers ...ResponseHandler) {
	r.before = append(r.before, handlers...)
}

// After appends handler(s) after all other handlers, globally for the instance of the router
func (r *Router) After(handlers ...ResponseHandler) {
	r.after = append(r.after, handlers...)
}

// Logger to use. Did works the same way as the other ResponseHandlers, except it does not exist
// in the middleware chain, and is being called even though next was never called
func (r *Router) Logger(logger ResponseHandler) {
	r.logger = logger
}

// Serve ..
func (r *Router) Serve(port ...int) {
	// Setting the default port, as we're using variadic variables, to make it possible to use Serve() parameterless
	portString := ":8080"
	if len(port) > 0 {
		portString = fmt.Sprintf(":%d", port[0])
	}

	fmt.Printf("Now serving on localhost%s\n", portString)

	// Binding the main request function to handle all requests made
	// This uses the default handleFunc, and makes all the logic from the same function
	http.HandleFunc("/", r.request)

	// Booting the server up, using the standard package http.
	log.Fatal(http.ListenAndServe(portString, nil))
}

// redirectPath fixes the path by either include a slash, or remove one.
// Searches for the fixed path and returns a boolean value for the result, and the redirect path
func (r *Router) redirectPath(path, method string) (bool, string) {
	redirectPath := path
	pathLength := len(redirectPath)
	if redirectPath[pathLength-1] == slash {
		redirectPath = redirectPath[:pathLength-1]
	} else {
		redirectPath += "/"
	}

	// Attempt to find the fixed route
	found, handler, _ := r.findRoute(redirectPath, method)

	return found && handler != nil, redirectPath
}

func (r *Router) request(w http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path

	found, handlers, params := r.findRoute(path, rq.Method)
	c := newContext()
	c.writer, c.request = w, rq

	if found && len(handlers) > 0 {
		handlerChain := []ResponseHandler{}

		if r.before != nil {
			handlerChain = append(handlerChain, r.before...)
		}

		handlerChain = append(handlerChain, handlers...)

		if r.after != nil {
			handlerChain = append(handlerChain, r.after...)
		}

		c.params, c.handlers, c.currentHandler, c.maxHandlers = Params{params}, handlerChain, 0, len(handlerChain)

		c.callByIndex(0)
	} else if found, redirectPath := r.redirectPath(path, rq.Method); found && r.RedirectSlashes {
		c.status = http.StatusMovedPermanently
		http.Redirect(w, rq, redirectPath, c.status)
	} else {
		c.status = http.StatusNotFound
		// Error handler here
		if r.NotFound == nil {
			fmt.Fprintln(w, "Requested page was not found")
		} else {
			r.NotFound(c)
		}

	}

	if r.logger != nil {
		r.logger(c)
	}

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
	// TODO - Make params object instead of map
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

			if position > pathLength || path[i:position] != res.path {
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

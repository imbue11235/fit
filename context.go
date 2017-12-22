package fit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context gets supplied to all functions conforming to the ResponseHandler, when the route get's inserted.
type Context struct {
	// Supplied handlers (middleware).
	// This also contains the response function to make it possible to call middleware before and after.
	handlers []ResponseHandler

	// Keeping track of the current handler to be able to call the next by using Next()
	currentHandler int

	// Keeping track of the amount of handlers to avoid calling len() over and over
	maxHandlers int

	// ResponseWriter supplied by http.HandleFunc
	writer http.ResponseWriter

	// Request supplied by http.HandleFunc
	request *http.Request

	// Current status set
	status int

	// Struct for handling the parameters set, if any.
	params Params

	// Shared for handling shared values between middleware(s)
	shared Shared
}

// ResponseHandler type, which every insertion of a route has to conform to, to receive the Context object
type ResponseHandler func(c *Context)

func newContext() *Context {
	return &Context{}
}

// Next calls the next middleware in the chain and returns a boolean value.
// The boolean value is determined by whether the new middleware was found or not.
func (c *Context) Next() bool {
	if c.currentHandler < c.maxHandlers-1 {
		c.currentHandler++
		c.callByIndex(c.currentHandler)

		return true
	}

	fmt.Println("There is no middleware left in the chain. Calling .Next() has no effect.")

	return false
}

// Writer returns an instance of the http.ResponseWriter.
func (c *Context) Writer() http.ResponseWriter {
	return c.writer
}

// Request returns an instance of the *http.Request.
func (c *Context) Request() *http.Request {
	return c.request
}

// callByIndex calls a handler by given index.
func (c *Context) callByIndex(index int) {
	c.handlers[index](c)
}

// setStatus sets the status of the header.
// It accepts multiple ints, to make it optional.
// The default status set, will be the http.StatusOk => 200
func (c *Context) setStatus(code ...int) {
	c.status = http.StatusOK
	if len(code) > 0 {
		c.status = code[0]
	}
	c.writer.WriteHeader(c.status)
}

// Status returns the current set status int.
func (c *Context) Status() int {
	return c.status
}

// JSON tries to encode the given interface as JSON using json.Marshal() and write the result to body, with the supplied status code
// Status code is optional and will use the default code set by the setStatus() function http.StatusOK => 200
// If it fails, it will set the status as http.StatusInternalServerError and output the error
func (c *Context) JSON(data interface{}, status ...int) {
	b, err := json.Marshal(data)

	if err != nil {
		c.writer.WriteHeader(http.StatusInternalServerError)
		c.writer.Write([]byte(err.Error()))
	}

	c.setStatus(status...)
	c.writer.Write(b)
}

package fit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// Middleware and handlers
	handlers       []ResponseHandler
	currentHandler int
	maxHandlers    int

	//
	writer  http.ResponseWriter
	request *http.Request

	status int

	//
	params Params

	// Shared

	shared map[string]interface{}
}

type ResponseHandler func(c *Context)

func newContext() *Context {
	return &Context{}
}

// Middleware
func (c *Context) Next() {
	if c.currentHandler < c.maxHandlers-1 {
		c.currentHandler++
		c.callByIndex(c.currentHandler)

		return
	}

	fmt.Println("There is no middleware left in the chain. Calling .Next() has no effect.")
}

// Original functions
func (c *Context) Writer() http.ResponseWriter {
	return c.writer
}

func (c *Context) Request() *http.Request {
	return c.request
}

// Calls handler by index
func (c *Context) callByIndex(index int) {
	c.handlers[index](c)
}

// Return functions
func (c *Context) setStatus(code []int) {
	c.status = http.StatusOK
	if len(code) > 0 {
		c.status = code[0]
	}
	c.writer.WriteHeader(c.status)
}

func (c *Context) JSON(data interface{}, status ...int) {
	b, err := json.Marshal(data)

	if err != nil {
		c.writer.WriteHeader(http.StatusInternalServerError)
		c.writer.Write([]byte(err.Error()))
	}

	c.setStatus(status)
	c.writer.Write(b)
}

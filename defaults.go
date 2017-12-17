package fit

import (
	"net/http"
)

// Default Not found

func notFoundHandler() ResponseHandler {

	return func(c *Context) {
		response := make(map[string]string)
		response["message"] = "The URL you requested was not found."
		c.JSON(response, http.StatusNotFound)
	}
}

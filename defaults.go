package fit

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// Default logger
func DefaultLogger() ResponseHandler {

	return func(c *Context) {
		statusColor := color.BgGreen

		if c.status == http.StatusNotFound {
			statusColor = color.BgYellow
		} else if c.status != http.StatusOK {
			statusColor = color.BgRed
		}

		col := color.New(statusColor)
		col.Printf("  %d  ", c.status)

		spacing := 20
		statusText := http.StatusText(c.status)
		restSpace := (spacing - len(statusText)) / 2
		spacingString := strings.Repeat(" ", restSpace)

		var preSpacing string
		if len(statusText)%2 != 0 {
			preSpacing = " "
		}
		fmt.Printf(" %s%s%s | %s \n", preSpacing+spacingString, statusText, spacingString, c.request.URL.Path)
	}

}

// Default Not found

func notFoundHandler() ResponseHandler {

	return func(c *Context) {
		response := make(map[string]string)
		response["message"] = "The URL you requested was not found."
		c.JSON(response, http.StatusNotFound)
	}
}

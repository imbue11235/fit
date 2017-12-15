package fit

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

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
		fmt.Printf(" %s - %s\n", http.StatusText(c.status), c.request.URL.Path)
	}

}

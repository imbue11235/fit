package main

import (
	fit "fit_router"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

type Message struct {
	Username string `json:"username"`
	APIToken string `json:"api_token"`
	Message  string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func DefaultLogger() fit.ResponseHandler {

	return func(c *fit.Context) {
		statusColor, status := color.BgGreen, c.Status()
		switch {
		case status >= 300 && status < 400:
			statusColor = color.BgHiBlue
		case status >= 400 && status < 500:
			statusColor = color.BgYellow
		case status >= 500:
			statusColor = color.BgRed
		}

		col := color.New(statusColor)
		col.Printf("  %d  ", status)

		spacing, statusText := 20, http.StatusText(status)
		restSpace := (spacing - len(statusText)) / 2
		spacingString := strings.Repeat(" ", restSpace)

		var preSpacing string
		if len(statusText)%2 != 0 {
			preSpacing = " "
		}
		fmt.Printf(" %s%s%s | %s \n", preSpacing+spacingString, statusText, spacingString, c.Request().URL.Path)
	}

}

// User - Example endpoint function
func User(c *fit.Context) {
	_, apiToken := c.Shared().Get("api_token")
	_, username := c.Parameters().GetByName("username")
	m := Message{username, apiToken.(string), fmt.Sprintf("You are allowed to view this page, because your name is '%s'.", username)}

	c.JSON(m)
}

// OnlyAllowUsersWithName - Example middleware
func OnlyAllowUsersWithName(username string) fit.ResponseHandler {
	return func(c *fit.Context) {
		_, value := c.Parameters().GetByName("username")
		if username != value {
			m := ErrorMessage{fmt.Sprintf("You are not allowed, your name needs to be '%s' to view this page.", username)}
			c.JSON(m, http.StatusUnauthorized)

			return
		}

		c.Shared().Set("api_token", "some_api_token")

		c.Next()
	}
}

func main() {
	router := fit.NewRouter()

	router.Logger = DefaultLogger()

	// http://localhost:<portString>/user/trump to view intended page
	// http://localhost:<portString>/user/somerandomname to view the middleware in effect
	router.Get("/", func(c *fit.Context) {
		c.JSON("Root")
	})

	router.Get("/user/:username/sa/:test", OnlyAllowUsersWithName("trump"), User).Where("username", "^[a-z]*$")

	router.Get("/test/route/:id", func(c *fit.Context) {
		_, value := c.Parameters().GetByName("id")
		c.JSON(Response{fmt.Sprintf("Id is %s and apikey is %s", value, c.Request().FormValue("apikey"))})
	})

	router.Get("/test/route-test/*something", func(c *fit.Context) {
		_, value := c.Parameters().GetByName("something")
		c.JSON(Response{fmt.Sprintf("Something is %s", value)})
	})

	router.Get("/broken-json", func(c *fit.Context) {
		broken := make(chan int)

		c.JSON(broken)
	})

	router.Serve(4000)
}

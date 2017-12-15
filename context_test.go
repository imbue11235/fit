package fit

import (
	"testing"
)

var (
	context = newContext()
)

func TestContextInitialization(t *testing.T) {
	if context == nil {
		t.Fatal("Context was never initialized")
	}
}

func HandlerTest1(t *testing.T) ResponseHandler {
	return func(c *Context) {
		c.Shared().Set("teststring", "Some string value")
		c.Shared().Set("testinteger", 5325)

		if c.shared["teststring"] != "Some string value" {
			t.Errorf("Shared value did not get set, expected '%s', got '%s'.", "Some string value", c.shared["teststring"])
		}

		if c.shared["testinteger"] != 5325 {
			t.Errorf("Shared value did not get set, expected '%d', got '%d'.", 5325, c.shared["testinteger"])
		}

		c.Next()
	}
}

func HandlerTest2(t *testing.T) ResponseHandler {
	return func(c *Context) {
		_, teststring := context.Shared().Get("teststring")
		if teststring != "Some string value" {
			t.Errorf("Shared value did not get passed to next function, expected '%s', got '%s'.", "teststring", teststring)
		}
		_, testinteger := context.Shared().Get("testinteger")
		if testinteger != 5325 {
			t.Errorf("Shared value did not get passed to next function, expected '%d', got '%d'.", 5325, testinteger)
		}
	}
}

func TestContextMiddleware(t *testing.T) {
	handlers := []ResponseHandler{HandlerTest1(t), HandlerTest2(t)}
	context.handlers = handlers
	context.maxHandlers = len(handlers)
	context.callByIndex(0)
}

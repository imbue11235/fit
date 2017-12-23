package fit

import (
	"net/http"
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

		if _, val := c.Shared().Get("teststring"); val != "Some string value" {
			t.Errorf("Shared value did not get set, expected '%s', got '%s'.", "Some string value", val)
		}

		if _, val := c.Shared().Get("testinteger"); val != 5325 {
			t.Errorf("Shared value did not get set, expected '%d', got '%d'.", 5325, val)
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

		// Should return false, as were calling a value that was never set
		ok, _ := context.Shared().Get("thisdoesnotexist")

		if ok {
			t.Errorf("Should not exist. Expected '%t', got '%t'", false, ok)
		}

		doesNextExist := c.Next()

		if doesNextExist {
			t.Errorf("There shouldn't be anymore middleware in the chain, but one was found")
		}
	}
}

func TestContextMiddlewareAndShared(t *testing.T) {
	handlers := []ResponseHandler{HandlerTest1(t), HandlerTest2(t)}
	context.handlers = handlers
	context.maxHandlers = len(handlers)
	context.callByIndex(0)
}

func TestContextParameters(t *testing.T) {
	testParams := make(map[string]string)
	testParams["id"] = "22"
	testParams["name"] = "John"
	context.params = Params{testParams}

	ok, id := context.Parameters().GetByName("id")

	if !ok || id != "22" {
		t.Errorf("Value for 'id' was wrong. Expected '%s', got '%s'", "22", id)
	}

	ok2, name := context.Parameters().GetByName("name")

	if !ok2 || name != "John" {
		t.Errorf("Value for 'name' was wrong. Expected '%s', got '%s'", "John", name)
	}

	exists, _ := context.Parameters().GetByName("doesnotexist")

	if exists {
		t.Error("Shouldn't be able to fetch value for 'doesnotexist'")
	}
}

func TestContextStatus(t *testing.T) {
	context.status = http.StatusTeapot

	if context.Status() != http.StatusTeapot {
		t.Errorf("Status is not correct, expected %d, got %d", http.StatusTeapot, context.Status())
	}
}

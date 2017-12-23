package fit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	router = NewRouter()
)

func TestRouterInitialization(t *testing.T) {
	if router == nil {
		t.Fatal("Router was never initialized")
	}
}

/* Testing the router print. Not an actual test
func TestRoutePrinting(t *testing.T) {

	router.addRoute("/", []string{"GET"})
	router.addRoute("/test", []string{"GET"})
	router.addRoute("/team", []string{"GET"})

	router.addRoute("/teammate", []string{"GET"})
	router.addRoute("/teammate", []string{"POST"})
	router.addRoute("/testable", []string{"GET"})
	router.addRoute("/testamonials/:id", []string{"GET"})
	router.addRoute("/testamonials/:id/something/:anotherid", []string{"GET"})
	router.addRoute("/testamonials/:id/somethang/:anotherid", []string{"GET"})

	router.PrintTree()
}
*/

func TestInsertion(t *testing.T) {
	path := "/some/testing/path"
	router.addRoute(path, nil)

	// Testing if the insertion was succesful. Child should just contain the full path
	if router.res.children[0].path != path {
		t.Errorf("First insertion went wrong, expected '%s', got '%s'.", path, router.res.children[0].path)
	}

	path2 := "/some/teaming/path"
	router.addRoute(path2, nil)

	// Testing the second insertion and if the split was done correctly
	if router.res.children[0].path != path[:8] {
		t.Errorf("Second insertion. First child is wrong, expected '%s', got '%s'.", path[:8], router.res.children[0].path)
	}

	if router.res.children[0].children[0].path != path2[8:] {
		t.Errorf("Second insertion. First child of first child is wrong expected '%s', got '%s'.", path2[8:], router.res.children[0].children[0].path)
	}

	if router.res.children[0].children[1].path != path[8:] {
		t.Errorf("Second insertion. Second child of first child is wrong expected '%s', got '%s'.", path[8:], router.res.children[0].children[0].path)
	}

	wildcardPath := "/testing/*all"
	router.addRoute(wildcardPath, nil)

	if router.res.children[0].children[1].path != wildcardPath[1:9] {
		t.Errorf("Wildcard insertion. Wildcard was not found. Expected '%s', got '%s'", wildcardPath[1:9], router.res.children[0].children[1].path)
	}
}

func TestFindingParameterizedRoute(t *testing.T) {
	router.addRoute("/find/:this/:withid", []string{"GET"})

	found, _, params := router.findRoute("/find/something/23", "GET")

	if !found {
		t.Errorf("Expected to get %t but got %t", true, found)
	} else {
		expectedParams := make(map[string]string)
		expectedParams["this"] = "something"
		expectedParams["withid"] = "23"
		if !reflect.DeepEqual(expectedParams, params) {
			t.Errorf("Expected to find params %s, found %s", expectedParams, params)
		}
	}
}

// Test router requests

type testMessage struct {
	Message    interface{} `json:"message"`
	Parameters []string    `json:"parameter"`
}

type testRoute struct {
	insertRoute          string
	visitRoute           string
	method               string
	parameterIdentifiers []string
	expectedStatus       int
	expectedMessage      interface{}
	expectedParameters   []string
}

func TestRoutes(t *testing.T) {

	routes := []testRoute{
		// Fetch - Static Routes
		testRoute{"/", "/", "GET", nil, http.StatusOK, "Root", nil},
		testRoute{"/photos", "/photos", "GET", nil, http.StatusOK, "You've found photos!", nil},

		// Fetch - Parameterized
		testRoute{"/photos/:id", "/photos/23", "GET", []string{"id"}, http.StatusOK, "Showing photo 23", []string{"23"}},

		// Fetch - Everything
		testRoute{"/photos/:id/comment/*all", "/photos/24/comment/asft4asf433", "GET", []string{"id", "all"}, http.StatusOK, "Showing a comment", []string{"24", "asft4asf433"}},

		// Redirects - Static
		testRoute{"", "/photos/", "GET", nil, http.StatusMovedPermanently, "You've found photos!", nil},

		// Redirects - Parameterized
		testRoute{"", "/photos/23/", "GET", []string{"id"}, http.StatusMovedPermanently, "Showing photo 23", []string{"23"}},
		testRoute{"/comments/:id/", "/comments/57", "GET", []string{"id"}, http.StatusMovedPermanently, "Comment #57", []string{"57"}},

		// 404 - Static
		testRoute{"", "/photoas/", "GET", nil, http.StatusNotFound, "The URL you've requested was not found.", nil},

		// 404 - Parameterized

		// 500 - Invalid JSON
		testRoute{"/invalid-json", "/invalid-json", "GET", nil, http.StatusInternalServerError, make(chan int), nil},
	}

	for _, route := range routes {
		if route.insertRoute != "" {
			router.addRoute(route.insertRoute, []string{route.method}, func(c *Context) {
				message := testMessage{route.expectedMessage, nil}
				if route.expectedParameters != nil {
					for _, param := range route.parameterIdentifiers {
						_, val := c.Parameters().GetByName(param)

						message.Parameters = append(message.Parameters, val)
					}
				}
				c.JSON(message)
			})
		}
		findRoute(route, t)
	}

}

func findRoute(route testRoute, t *testing.T) {
	req := httptest.NewRequest(route.method, route.visitRoute, nil)
	w := httptest.NewRecorder()
	router.request(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != route.expectedStatus {
		t.Errorf("Status for '%s' code is wrong. Expected %d, got %d", route.visitRoute, route.expectedStatus, resp.StatusCode)
	}

	var m testMessage
	err := json.Unmarshal(body, &m)
	if err == nil && m.Message != route.expectedMessage {
		t.Errorf("Message in body is wrong for '%s'. Expected '%s', got '%s'", route.visitRoute, route.expectedMessage, m.Message)
	}

	if route.expectedParameters != nil && err == nil {
		if !reflect.DeepEqual(m.Parameters, route.expectedParameters) {
			t.Errorf("Parameter in body is wrong for '%s'. Expected '%s', got '%s'", route.visitRoute, route.expectedParameters, m.Parameters)
		}
	}

}

/*
func TestRouter404(t *testing.T) {
	req := httptest.NewRequest("GET", "/a/route/myarticle/asgdsgds", nil)
	w := httptest.NewRecorder()
	router.request(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Status code is wrong. Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

}

func TestRouterLogger(t *testing.T) {
	//router.Logger =
}
*/

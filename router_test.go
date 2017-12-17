package fit

import (
	"encoding/json"
	"fmt"
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

func SomeHandler() {
	fmt.Println("Works")
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

type testMessage struct {
	Message string `json:"message"`
	Article string `json:"article"`
}

func TestSetupRouterRequest(t *testing.T) {

	router.addRoute("/a/route/:article", []string{"GET"}, func(c *Context) {
		_, value := c.Parameters().GetByName("article")

		message := testMessage{"Hey, it worked!", value}
		c.JSON(message)
	})

}

func TestRouterRequest(t *testing.T) {

	req := httptest.NewRequest("GET", "/a/route/myarticle", nil)
	w := httptest.NewRecorder()
	router.request(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var m testMessage
	err := json.Unmarshal(body, &m)
	if err == nil && m.Message != "Hey, it worked!" {
		t.Errorf("Body is wrong. Expected '%s', got '%s'", "Hey, it worked!", m.Message)
	}

	if err == nil && m.Article != "myarticle" {
		t.Errorf("Body is wrong. Expected '%s', got '%s'", "myarticle", m.Article)
	}
}

func TestRouterRequestExtraSlash(t *testing.T) {

	router.addRoute("/a/testing-route/:article", []string{"GET"}, func(c *Context) {})

	req := httptest.NewRequest("GET", "/a/route/myarticle/", nil)
	w := httptest.NewRecorder()
	router.request(w, req)

	resp := w.Result()

	// If the status code is 301, we found the route, even though it had an extra slash ("/"), and are getting redirected.
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Status code is wrong. Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

}

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

package fit

import (
	"fmt"
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

/*
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
}

func TestFindingParameterizedRoute(t *testing.T) {
	router.addRoute("/find/:this/:withid", []string{"GET"})

	found, _, params := router.Get("/find/something/23", "GET")

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
*/

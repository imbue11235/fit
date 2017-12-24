package fit

import (
	"fmt"
	"testing"
)

type benchmarkRoute struct {
	method string
	path   string
}

var (
	benchmarkRouter = NewRouter()
	// https://developers.google.com/+/web/api/rest/latest/
	googlePlusApi = []benchmarkRoute{
		{"GET", "/people/:userId"},
		{"GET", "/people"},
		{"GET", "/activities/:activityId/people/:collection"},
		{"GET", "/people/:userId/people/:collection"},
		{"GET", "/people/:userId/openIdConnect"},
		{"GET", "/people/:userId/activities/:collection"},
		{"GET", "/activities/:activityId"},
		{"GET", "/activities"},
		{"GET", "/activities/:activityId/comments"},
		{"GET", "/activities/:activityId/comments/:comment/user/:userId/favorites/:favoriteId/latest/:lastestId"},
		{"GET", "/comments/:commentId"},
		{"POST", "/people/:userId/moments/:collection"},
		{"GET", "/people/:userId/moments/:collection"},
		{"DELETE", "/moments/:id"},
		{"GET", "/custom/*all"},
	}
)

func TestMain(m *testing.M) {
	insertRoutesForTesting()
	m.Run()
}

func insertRoutesForTesting() {
	for _, route := range googlePlusApi {
		benchmarkRouter.addRoute(route.path, []string{route.method}, func(c *Context) {
			fmt.Println(route.path)
		})
	}
}

func benchmarkFind(path, method string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkRouter.findRoute(path, method)
	}
}

func BenchmarkStaticRoute(b *testing.B) {
	benchmarkFind("/people", "GET", b)
}

func Benchmark1ParameterRoute(b *testing.B) {
	benchmarkFind("/activities/44", "GET", b)
}

func Benchmark2ParameterRoute(b *testing.B) {
	benchmarkFind("/people/22/moments/5", "GET", b)
}

func Benchmark5ParameterRoute(b *testing.B) {
	benchmarkFind("/activities/22/comments/1/user/45645/favorites/242/latest/435", "GET", b)
}

func BenchmarkCatchAllRoute(b *testing.B) {
	benchmarkFind("/custom/some-custom-string", "GET", b)
}

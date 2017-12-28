package fit

import (
	"testing"
)

type benchmarkRoute struct {
	method string
	path   string
	visit  string
}

var (
	benchmarkRouter = NewRouter()
	// https://developers.google.com/+/web/api/rest/latest/
	googlePlusAPI = []benchmarkRoute{
		{"GET", "/people/:userId", "/people/23"},
		{"GET", "/people", "/people"},
		{"GET", "/activities/:activityId/people/:collection", "/activities/100/people/66"},
		{"GET", "/people/:userId/people/:collection", "/people/65/people/323"},
		{"GET", "/people/:userId/openIdConnect", "/people/235/openIdConnect"},
		{"GET", "/people/:userId/activities/:collection", "/people/12/activities/34657"},
		{"GET", "/activities/:activityId", "/activities/2346"},
		{"GET", "/activities", "/activities"},
		{"GET", "/activities/:activityId/comments", "/activities/346/comments"},
		{"GET", "/activities/:activityId/comments/:comment/user/:userId/favorites/:favoriteId/latest/:lastestId", "/activities/22/comments/1/user/45645/favorites/242/latest/435"},
		{"GET", "/comments/:commentId", "/comments/235"},
		{"POST", "/people/:userId/moments/:collection", "/people/234657/moments/23543"},
		{"GET", "/people/:userId/moments/:collection", "/people/324675/moments/332"},
		{"DELETE", "/moments/:id", "/moments/23"},
		{"GET", "/custom/*all", "/custom/whatever-string"},
	}
)

func TestMain(m *testing.M) {
	insertRoutesForTesting()
	m.Run()
}

func insertRoutesForTesting() {
	for _, route := range googlePlusAPI {
		benchmarkRouter.addRoute(route.path, []string{route.method}, nil)
	}
}

func benchmarkFind(path, method string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkRouter.findRoute(path, method)
	}
}

func BenchmarkFindStaticRoute(b *testing.B) {
	benchmarkFind("/people", "GET", b)
}

func BenchmarkFind1ParameterRoute(b *testing.B) {
	benchmarkFind("/activities/44", "GET", b)
}

func BenchmarkFind2ParameterRoute(b *testing.B) {
	benchmarkFind("/people/22/moments/5", "GET", b)
}

func BenchmarkFind5ParameterRoute(b *testing.B) {
	benchmarkFind("/activities/22/comments/1/user/45645/favorites/242/latest/435", "GET", b)
}

func BenchmarkFindCatchAllRoute(b *testing.B) {
	benchmarkFind("/custom/some-custom-string", "GET", b)
}

func BenchmarkFindAllRoutes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, route := range googlePlusAPI {
			benchmarkRouter.findRoute(route.visit, route.method)
		}
	}
}

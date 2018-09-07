# Fit router

A simple router implementation based on a radix trie

## In development

### Usage

Simple route binding

```go
router := fit.NewRouter()

router.Get("/", func(c *fit.Context) {
    message := struct {
        Message: string `json:"message"`
    } {
        "Hi there"
    }

    c.JSON(message)
})

// router.Post() etc..

router.Serve(3000) // Serve on port 3000
```

Middleware

```go
// Can be set globally
router.Before(func(c *fit.Context){ /* code */ }) // Triggers before requests
router.After(func(c *fit.Context){ /* code */ }) // Triggers after requests

// Or route specific
router.Get("/", someMiddlewareFunc, myActualEndpointFunc, someHandlerAfterRequestFunc)

```

More examples are coming

### Benchmarks

```go
BenchmarkFindStaticRoute-8       	50000000	        26.0 ns/op
BenchmarkFind1ParameterRoute-8   	20000000	       117 ns/op
BenchmarkFind2ParameterRoute-8   	10000000	       147 ns/op
BenchmarkFind5ParameterRoute-8   	10000000	       196 ns/op
BenchmarkFindCatchAllRoute-8     	20000000	       115 ns/op
BenchmarkFindAllRoutes-8         	 1000000	      2016 ns/op
```

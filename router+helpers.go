package fit

import (
	"fmt"
	"strings"
)

// Route helpers
func (r *Router) Get(path string, handlers ...ResponseHandler) *Options {
	return r.addRoute(path, []string{"GET"}, handlers...)
}

// Route find

func find(src string, target byte, start int, pathLength int) int {
	i := start

	for ; i < pathLength && src[i] != target; i++ {
	}
	return i
}

// Route print tree

func (r *Router) PrintTree() {
	res := r.res
	printResource(res, 0, false)
}

func printResource(res *resource, amount int, pre bool) {
	total := len(res.children)
	for i, r := range res.children {
		var methods string
		if len(r.methods) > 0 {
			keys := make([]string, 0, len(r.methods))
			for k := range r.methods {
				keys = append(keys, k)
			}
			methods = fmt.Sprint("(", strings.Join(keys, "|"), ")")
		}

		flag := "└──"

		if i == 0 && total > i+1 {
			flag = "├──"
		}

		var spacing string
		if amount != 0 {
			spacing = strings.Repeat("     ", amount-1)
		}

		if pre {
			spacing += "|   "
		} else {
			spacing += "    "
		}

		fmt.Println(spacing, flag, r.path, methods)

		if len(r.children) > 0 {
			printResource(r, amount+1, total != i+1)
		}
	}
}

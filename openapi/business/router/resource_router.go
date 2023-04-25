package router

import (
	"strings"
)

type ResourceRouter struct {
	trees map[string]*node
}

func NewResourceRouter() *ResourceRouter {
	rr := &ResourceRouter{}

	rr.trees = make(map[string]*node)

	return rr
}

func (rr *ResourceRouter) Add(method, path string) {
	method = strings.ToUpper(method)
	var root *node
	root, ok := rr.trees[method]
	if !ok {
		root = &node{
			root: true,
			Path: "/",
		}
	}
	root.add(path, nil, "")
	rr.trees[method] = root
}

func (rr *ResourceRouter) Find(method, path string) bool {
	if rr.trees == nil {
		return false
	}
	if node, ok := rr.trees[method]; ok {
		pass, _, _ := node.find(path)
		return pass
	}
	return false
}

func (rr *ResourceRouter) IsEmpty() bool {
	return len(rr.trees) == 0
}

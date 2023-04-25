package router

import (
	"blogrpc/openapi/business/controller"
)

var Trees map[string]*node

func init() {
	Trees = make(map[string]*node)
}

func Handle(method, path string, action *controller.ControllerAction) {
	var root *node
	root, ok := Trees[method]
	if !ok {
		root = &node{
			root: true,
			Path: "/",
		}
	}
	root.add(path, action, "")
	Trees[method] = root
}

func AddByProxy(method, path string, proxyPath string) {
	var root *node
	root, ok := Trees[method]
	if !ok {
		root = &node{
			root: true,
			Path: "/",
		}
	}
	root.add(path, nil, proxyPath)
	Trees[method] = root
}

func Find(method, path string) (find bool, params map[string]string, node *node) {
	if Trees == nil {
		return
	}
	if node, ok := Trees[method]; ok {
		return node.find(path)
	}
	return
}

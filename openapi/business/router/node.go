package router

import (
	"strings"

	"blogrpc/openapi/business/controller"
)

type node struct {
	root     bool
	wild     bool
	children []*node

	Path      string
	Action    *controller.ControllerAction
	ProxyPath string
}

func (n *node) add(path string, action *controller.ControllerAction, proxyPath string) {
	path = strings.Trim(path, "/")
	if path == "" {
		n.Action = action
		return
	}
	paths := strings.Split(path, "/")
	n.addRoute(paths, action, proxyPath)
	return
}

func (n *node) addRoute(paths []string, action *controller.ControllerAction, proxyPath string) {
	path := paths[0]
	wild := isWild(path)
	children := n.children
	depth := len(paths)
	if len(children) < 1 {
		// No child
		child := &node{
			Path:     path,
			root:     false,
			wild:     wild,
			children: nil,
			Action:   nil,
		}
		n.children = []*node{child}
		if depth == 1 {
			child.Action = action
			child.ProxyPath = proxyPath
			return
		}
		child.addRoute(paths[1:], action, proxyPath)
		return
	}
	for _, child := range children {
		if child.Path == path {
			if depth > 1 {
				child.addRoute(paths[1:], action, proxyPath)
				return
			}
			if child.Action == nil {
				child.Action = action
			}

			child.ProxyPath = proxyPath
			return
		}
	}
	// No matched child
	child := &node{
		Path:     path,
		root:     false,
		wild:     wild,
		children: nil,
		Action:   nil,
	}
	n.children = append(n.children, child)
	if len(paths) == 1 {
		child.Action = action
		child.ProxyPath = proxyPath
		return
	}
	child.addRoute(paths[1:], action, proxyPath)
	return
}

func (n *node) Find(path string) (pass bool, params map[string]string, node *node) {
	return n.find(path)
}

func (n *node) find(path string) (pass bool, params map[string]string, node *node) {
	path = strings.Trim(path, "/")
	if path == "" {
		// root action
		node = n
		pass = true
		return
	}

	paths := strings.Split(path, "/")
	params = make(map[string]string)
	pass, node = n.get(paths, params)
	return
}

func (n *node) get(paths []string, params map[string]string) (pass bool, currentNode *node) {
	children := n.children
	if len(children) == 0 {
		return
	}
	depth := len(paths)
	path := paths[0]
	var wildChild *node
	for _, child := range children {
		if child.wild {
			wildChild = child
			continue
		}
		if child.Path == path {
			if depth == 1 {
				pass = true
				currentNode = child
				return
			}
			return child.get(paths[1:], params)
		}
	}

	if wildChild == nil {
		return
	}

	if params == nil {
		params = make(map[string]string)
	}
	params[removeWildcard(wildChild.Path)] = path
	if depth == 1 {
		pass = true
		currentNode = wildChild
		return
	}

	return wildChild.get(paths[1:], params)
}

func isWild(path string) bool {
	c := path[0]
	return c == ':' || c == '*'
}

func removeWildcard(path string) string {
	path = strings.TrimPrefix(path, "*")
	return strings.TrimPrefix(path, ":")
}

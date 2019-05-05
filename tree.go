// Copyright 2019 The liuyousama/Router Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"net/http"
	"strings"
)

type Tree struct {
	root     *Node
	children map[string]*Node
}

type Node struct {
	children    map[string]*Node
	paramPath   string
	handle      http.HandlerFunc
	middlewares []middlewareFunc
}

//NewTree returns a new Tree.
func NewTree() *Tree {
	return &Tree{root: NewNode(), children: make(map[string]*Node)}
}

//NewNode returns a new Node.
func NewNode() *Node {
	return &Node{children: make(map[string]*Node)}
}

//AddPath add the custom path to the tree.
func (tree *Tree) AddPath(path string, r *Router, handler http.HandlerFunc) {
	//if path is empty,handler will be on the root node.
	if path == "" {
		tree.root.handle = handler
		tree.root.middlewares = append(tree.root.middlewares, r.middlewares...)
		return
	}

	pathList := strings.Split(path, "/")
	node := tree.root

	//foreach every node
	for _, path = range pathList {
		//if a node is with a param,check it weather right.
		if strings.HasPrefix(path, ":") {
			if node.paramPath != "" {
				panic(fmt.Errorf("Router Conflict!!"))
			}
			node.paramPath = path
			newNode := NewNode()
			node.children[path] = newNode
			node = newNode
			continue
		}

		newNode, ok := node.children[path]
		if !ok {
			newNode = NewNode()
			node.children[path] = newNode
		}
		node = newNode
	}

	node.handle = handler
	node.middlewares = append(node.middlewares, r.middlewares...)
}

//Find finds the node with a pathList and record the param in the pathList
func (t *Tree) Find(pathList []string) (*Node, map[string]string) {
	node := t.root
	paramMap := make(map[string]string)

	for _, path := range pathList {
		newNode, ok := node.children[path]

		//if there is not the next node matched with the path and current node has a param child node, record the param.
		if !ok && node.paramPath != "" {
			key := strings.TrimPrefix(node.paramPath, ":")
			paramMap[key] = path
			node = node.children[node.paramPath]
			continue
		}

		if !ok {
			return nil, nil
		}
		node = newNode
	}
	return node, paramMap
}

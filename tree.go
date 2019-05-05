package router

import (
	"strings"
	"net/http"
	"fmt"
)

type Tree struct {
	root *Node
	children map[string]*Node
}

type Node struct {
	children map[string]*Node
	paramPath string
	handle http.HandlerFunc
	middlewares []middlewareFunc
}

func NewTree() *Tree{
	return &Tree{root:NewNode(),children:make(map[string]*Node)}
}

func NewNode() *Node {
	return &Node{children:make(map[string]*Node)}
}

//将给定的path与handler添加到指定的Router中
func (tree *Tree) AddPath(path string,r *Router,handler http.HandlerFunc) {
	//如果传入的path为空，就直接将handler添加到根结点上
	if path == "" {
		tree.root.handle = handler
		tree.root.middlewares = append(tree.root.middlewares, r.middlewares...)
		return
	}
	
	pathList := strings.Split(path, "/")
	node := tree.root
	
	//循环每一个路由结点
	for _,path = range pathList {
		//如果这个节点为参数形式(:name)，单独处理
		if strings.HasPrefix(path,":") {
			if node.paramPath!="" {
				panic(fmt.Errorf("router confilct!!"))
			}
			node.paramPath = path
			newNode := NewNode()
			node.children[path] = newNode
			node = newNode
			continue
		}
		
		newNode,ok := node.children[path];
		if !ok{
			newNode = NewNode()
			node.children[path] = newNode
		}
		node = newNode
	}
	
	node.handle = handler
	node.middlewares = append(node.middlewares, r.middlewares...)
}

func (t *Tree) Find(pathList []string) *Node {
	node := t.root
	for _, path := range pathList  {
		newNode, ok := node.children[path];
		if !ok{
			return nil
		}
		node = newNode
	}
	return node
}

// Copyright 2019 The liuyosama/Router Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

//Router:定义路由容器对象结构体，使用树形结构存储路由的不同节点
type Router struct {
	group        string
	trees        map[string]*Tree
	middlewares  []middlewareFunc
	notFoundPage http.HandlerFunc
}

type (
	//middleware方法签名
	middlewareFunc func(next http.HandlerFunc) http.HandlerFunc
	//用来存储参数的上下文key
	contextKeyType struct{}
)

var contextKey = contextKeyType{}

//New:New方法返回一个默认的路由对象
func New() *Router {
	return &Router{trees: make(map[string]*Tree)}
}

//Group:Group方法返回一个分组的路由对象
func (r *Router) Group(groupPath string) *Router {
	return &Router{
		trees: r.trees,
		group: groupPath}
}

//GET:GET方法添加一个get请求的path
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}
//POST:POST方法添加一个get请求的path
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}
//DELETE:DELETE方法添加一个get请求的path
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}
//PUT:PUT方法添加一个get请求的path
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}
//PATCH:PATCH方法添加一个get请求的path
func (r *Router) PATCH(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handler)
}

//Handler:Handler处理新增的路由，将其放入路由树结构中
func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if method != http.MethodGet &&
		method != http.MethodPost &&
		method != http.MethodDelete &&
		method != http.MethodPut &&
		method != http.MethodPatch {
		panic(fmt.Errorf("invaild method!"))
	}

	tree, ok := r.trees[method]
	if !ok {
		tree = NewTree()
		r.trees[method] = tree
	}

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	if r.group != "" {
		path = r.group + "/" + path
	}

	tree.AddPath(path, r, handler)

}

//NotFoundPage:为router设置404handler
func (r *Router) NotFoundPage(handler http.HandlerFunc) {
	r.notFoundPage = handler
}

//Use:Use为当前路由对象设置middleware
func (r *Router) Use(middlewares ...middlewareFunc) {
	if len(middlewares) > 0 {
		r.middlewares = append(r.middlewares, middlewares...)
	}
}
//ServeHTTP:ServeHTTP实现路由签名
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//检查是否有用户自定义的404handler，如果没有，则使用默认的。并且在函数末尾调用404handler
	if r.notFoundPage == nil {
		r.notFoundPage = defaultNotFoundPage
	}
	defer r.notFoundPage(w, req)

	//判断方法对应的tree是否存在
	tree, ok := r.trees[req.Method]
	if ok {
		//获取请求路径，转化为路径节点列表
		path := req.URL.Path
		path = strings.TrimPrefix(path, "/")
		pathList := strings.Split(path, "/")

		//如果路径节点为0，就执行根结点对应的handler,否则在树中查找是否有对应节点
		if len(pathList) == 0 {
			if tree.root.handle != nil {
				handle(w, req, tree.root.handle, tree.root.middlewares)
			}
		} else if len(pathList) > 0 {
			node, paramMap := tree.Find(pathList)
			//将参数列表存入request上下文中去
			ctx := context.WithValue(req.Context(), contextKey, paramMap)
			req = req.WithContext(ctx)

			if node != nil && node.handle != nil {
				handle(w, req, node.handle, node.middlewares)
			}
		}

	}
}

//GetAllParams：GetAllParams获取当前请求的所有参数
func GetAllParams(r *http.Request) map[string]string {
	paramMap, ok := r.Context().Value(contextKey).(map[string]string)
	if ok {
		return paramMap
	} else {
		return nil
	}
}

//GetParam：GetParam获取当前请求特定key的参数
func GetParam(r *http.Request, key string) string {
	val, ok := GetAllParams(r)[key]
	if ok {
		return val
	} else {
		return ""
	}
}

//handle:handle处理handler与middleware
func handle(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, middlewares []middlewareFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	handler(w, r)
}
//defaultNotFoundPage:defaultNotFoundPage为默认的404Handler
func defaultNotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404! Not Found Page!!")
}

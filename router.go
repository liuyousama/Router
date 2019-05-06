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

//Router is a container of route.
type Router struct {
	group        string
	trees        map[string]*Tree
	middlewares  []middlewareFunc
	notFoundPage http.HandlerFunc
}

type (
	//middleware func signature.
	middlewareFunc func(next http.HandlerFunc) http.HandlerFunc
	//a context ket used to store the params.
	contextKeyType struct{}
)

var contextKey = contextKeyType{}

//New returns a default router.
func New() *Router {
	return &Router{trees: make(map[string]*Tree)}
}

//Group returns a group router.
func (r *Router) Group(groupPath string) *Router {
	return &Router{
		trees: r.trees,
		group: groupPath}
}

//GET add a route with GET method.
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

//POST add a route with POST method.
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

//DELETE add a route with DELETE method.
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

//PUT add a route with PUT method.
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

//PATCH:PATCH add a route with PATCH method.
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

//NotFoundPage set custom 404 handler to a router
func (r *Router) NotFoundPage(handler http.HandlerFunc) {
	r.notFoundPage = handler
}

//Use add middlewares to current router.
func (r *Router) Use(middlewares ...middlewareFunc) {
	if len(middlewares) > 0 {
		r.middlewares = append(r.middlewares, middlewares...)
	}
}

//ServeHTTP implements a router signature.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//check whether the custom 404 handler exists,if not,use the default 404 handler.
	if r.notFoundPage == nil {
		r.notFoundPage = defaultNotFoundPage
	}
	//call the 404 handler in the end
	defer r.notFoundPage(w, req)

	//check whether the tree with request method exists.
	tree, ok := r.trees[req.Method]
	if ok {
		//get the request path, and transform it into a list.
		path := req.URL.Path
		path = strings.TrimPrefix(path, "/")
		pathList := strings.Split(path, "/")

		//if the length of pathList is 0, check the root node handler.
		if len(pathList) == 0 {
			if tree.root.handle != nil {
				handle(w, req, tree.root.handle, tree.root.middlewares)
				return
			}
		} else if len(pathList) > 0 {
			node, paramMap := tree.Find(pathList)
			//store the param map into the request context.
			ctx := context.WithValue(req.Context(), contextKey, paramMap)
			req = req.WithContext(ctx)

			if node != nil && node.handle != nil {
				handle(w, req, node.handle, node.middlewares)
				return
			}
		}

	}
}

//GetAllParams gets the map with all params.
func GetAllParams(r *http.Request) map[string]string {
	paramMap, ok := r.Context().Value(contextKey).(map[string]string)
	if ok {
		return paramMap
	} else {
		return nil
	}
}

//GetParam gets the param with custom key.
func GetParam(r *http.Request, key string) string {
	val, ok := GetAllParams(r)[key]
	if ok {
		return val
	} else {
		return ""
	}
}

//handle handle the handler and middleware.
func handle(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, middlewares []middlewareFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	handler(w, r)
}

//defaultNotFoundPage is the default handler when 404 not fount is happen.
func defaultNotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404! Not Found Page!!")
}

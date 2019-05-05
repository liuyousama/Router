# Router    <a href="https://travis-ci.org/liuyousama/Router"><img src="https://travis-ci.com/liuyousama/Router.svg?branch=master" alt="Build Status"></a>  [![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/liuyousama/Router/master/LICENSE)    [![Release](https://img.shields.io/badge/release-v1.0-blue.svg?style=flat-square)](https://github.com/liuyousama/Router/releases/tag/v1.0)
`liuyousama/Router`是为Golang设计的的一个简洁，快速的路由框架，可以用来快速构建restful服务（A simple router framework for Golang）
## 安装（Install）
```
go get github.com:liuyousama/Router
```

## 特性（Featurs）
- RESTFUL风格
- 路由参数
- 路由分组
- 中间件
- 自定义 404 NotFound Handler

## 快速开始（Quick Start）
```go
func main() {
    r := Router.New()
    
    r.GET("/helloRouter",HelloRouterHandler)
    
    http.ListenAndServe(":8080", r)
}

func HelloRouterHander(w http.ResponseWriter, r *http.Request)  {
    w.Write([]Byte("Hello!Router!"))
}
```

## RESTFUL风格
```go
r := Router.New()
    
r.GET("/user",GetHandler)
r.POST("/user",PostHandler)
r.PUT("/user",PutHandler)
r.DELETE("/user",DeleteHandler)
r.PATCH("/user",PatchHandler)
```
## 路由参数
- 设置参数
```go
r.GET("/user/:id",GetHandler)
```
- 获取参数
```go
func GetHandler(w http.ResponseWriter, r *http.Request) {
	var paramMap map[string]string
    paramMap = Router.GetAllParams(r)
    id := paramMap["id"]
    //或者
    id := Router.GetParam(r, "id")
}
```
## 路由分组
```go
r := Router.New()

userGroup := r.Group("user")
{
	userGroup.GET("/list", GetUserListHandler)
}
```

## 使用中间件
- 使用全局路由
```go

r := Router.New()
r.Use(Logging,WriteHeader)

func Logging(next http.HanderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    		log.Printf("New request comes from %s",r.RemoteAddr))
    		next.ServeHTTP(w, r)
    	}
}

func Logging(next http.HanderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		    w.WriteHeader(http.StatusOk)
    		next.ServeHTTP(w, r)
    	}
}
```


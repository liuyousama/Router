# Router
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
    
    r.GET("helloRouter",HelloRouterHandler)
    
    http.ListenAndServe(":8080", r)
}

func HelloRouterHander(w http.ResponseWriter, r *http.Request)  {
    w.Write([]Byte("Hello!Router!"))
}
```

##RESTFUL风格

##路由参数

##路由分组

##使用中间件


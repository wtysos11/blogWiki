# 官方net/http包源码阅读

标签：golang http

## 参考资料

* [golang http serve源码阅读](https://www.ctolib.com/topics-3437.html)
* [go语言标准库阅读](https://books.studygolang.com/The-Golang-Standard-Library-by-Example/)
* [github:golang/go/src/net](https://github.com/golang/go/tree/master/src/net)

## 官方文档阅读

[http包文档](https://golang.org/pkg/net/http/)

http包提供了http客户端和服务端的实现。

提供了Get,Head,Post和PostForm来进行http或https的请求：

```golang
resp, err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
	url.Values{"key": {"Value"}, "id": {"123"}})
```

客户端在最后必须关闭response Body。

```golang
resp, err := http.Get("http://example.com/")
if err != nil {
	// handle error
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
// ...
```

如果想要操控http客户端头，重定向政策或其他设定，需要创建一个Client：

```golang
client := &http.Client{
	CheckRedirect: redirectPolicyFunc,
}

resp, err := client.Get("http://example.com")
// ...

req, err := http.NewRequest("GET", "http://example.com", nil)
// ...
req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req)
// ...
```

使用Transport来操纵代理、TLS设置、keep-alives、亚索其他设定。

```golang
tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```

Clients和Transport是并行安全的，可以同时被多个goroutine所使用而不用担心出现数据冲突。而且为了效率，这两个都应该只被创建一次。

ListenAndServe方法开启一个http服务器，使用给定的地址和处理句柄。句柄通常为nil，代表着使用默认的路由DefaultServeMux。使用Handle和HandleFunc方法向DefaultServeMux添加句柄：

```golang
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))
```

通过Server结构体，更多对服务器的行为是可以被实现的：

```golang
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
```

### 重要方法

* `func Handle(pattern string, handler Handler)`,Handle方法使用给定的模式向DefaultServeMux注册handler。ServeMux的文档详细描述了这个模式。
* `func HandleFunc(pattern string, handler func(ResponseWriter, *Request))`，HandleFunc和上文一样，只不过使用的函数类型不同。
* `func ListenAndServe(addr string, handler Handler) error`，启动一个TCP网络连接
* `func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error`，启动一个https网络连接。
* `func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error`，接受到来的https连接请求，对于每个请求创建一个新的service线程。线程读取请求，并调用handler来进行处理。

### 重要对象

* [Client](https://golang.org/pkg/net/http/#Client)
* [Cookie](https://golang.org/pkg/net/http/#Cookie)
* [Handler](https://golang.org/pkg/net/http/#Handler)
* [Header](https://golang.org/pkg/net/http/#Header)
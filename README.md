Pigeon
===

An easy-to-use, flexible HTTP web framework based on [Gin][gin], 
inspired by [nginx][nginx], [openresty][openresty] and [beego][beego].

Quick Start
---

### create file `myserver.go`

```go
package main

import (
	"github.com/opencurve/pigeon"
)

func main() {
	server := pigeon.NewHTTPServer()
	server.Route("/", func(r *pigeon.Request) bool {
		return r.Exit(200, "Hello Pigeon\n")
	})
	
	pigeon.Serve(server)
}
```

### build and run
```shell
$ go build myserver.go
$ ./myserver start
```

### access web server in default address
```shell
$ curl 127.0.0.1:8000
Hello Pigeon
```

[gin]: https://github.com/gin-gonic/gin
[nginx]: https://github.com/nginx/nginx
[openresty]: https://github.com/openresty
[beego]: https://github.com/beego/beego
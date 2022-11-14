# Ninja

[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://github.com/infinitasx/easi-go-aws)

> a simple ioc for https://github.com/facebookgo/inject

## ⚙ Installation

```text
go get -u github.com/eininst/ninja
```

## ⚡ Quickstart

```go
package main

import (
	"github.com/eininst/ninja"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
    Cfg    *ninja.Config `inject:""`
}

func (r *Router) Init(router fiber.Router) {
    router.Get("/", func(ctx *fiber.Ctx) error {
        return ctx.SendString("helloword")
    })
}

func main() {
    nj := ninja.NewApp("./examples/helloword.yml")
    ninja.Provide(new(Router))
	
    nj.Listen()
}
```


## Config yaml
```yaml
ninja:
  port: 8001
  middleware:
    - recover:
    - limiter:
        rate: 10
        burst: 200
        timeout: 5
    - status:
    - logger:
    - metrics:

---
profile: dev

redis:
  addr: kubernetes.docker.internal:6379
  db: 0
  poolSize: 100
  minIdleCount: 20

mysql:
  dsn: nft:Aa505814@tcp(localhost:3306)/credit?charset=utf8mb4&parseTime=True&loc=Local
  maxIdleCount: 32
  maxOpenCount: 128
  maxLifetime: 7200

rs:
  prefix: MQ_
  sender:
    maxLen: 100

---
profile: prod

redis:
  addr: kubernetes.docker.internal:6379
  db: 0
  poolSize: 100
  minIdleCount: 20

mysql:
  dsn: test:test123@tcp(localhost:3306)/credit?charset=utf8mb4&parseTime=True&loc=Local
  maxIdleCount: 32
  maxOpenCount: 128
  maxLifetime: 7200
```


## Depends
https://github.com/gofiber/fiber #Fiber framework

https://github.com/facebookgo/inject #依赖注入

https://github.com/eininst/fiber-prefork-grace #Fiber Prefork(多进程)模式下的优雅退出

https://github.com/eininst/fiber-middleware-burst #Fiber的令牌桶限流器

https://github.com/eininst/fiber-middleware-recover #Fiber的错误处理中间件

https://github.com/eininst/fiber-middleware-redoc #Fiber redoc


> See [examples](/examples)
## License

*MIT*
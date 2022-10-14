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
    Router fiber.Router  `inject:""`
    Cfg    *ninja.Config `inject:""`
}

func (r *Router) Init() {
    r.Router.Get("/", func(ctx *fiber.Ctx) error {
        return ctx.SendString("helloword")
    })
}

func main() {
    nj := ninja.New("./examples/helloword.yml")
    nj.Install(&Router{})

    nj.Listen()
}

```
> See [examples](/examples)
## License

*MIT*
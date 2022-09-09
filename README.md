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
	"fastgo/api/helloword"
	"github.com/eininst/ninja"
	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    var hapi helloword.Api
	
    ninja.Install(&hapi, app)

    app.Listen(":8080")
}
```

> See helloword.Api

```go
package helloword

import (
    v1 "fastgo/api/helloword/v1"
    "github.com/gofiber/fiber/v2"
)

type Api struct {
    App       *fiber.App       `inject:""`
    Helloword *v1.HellowordApi `inject:""`
}

func (api *Api) Init() {
    api.App.Get("/add", api.Helloword.Add)
}
```

```go
package v1

import (
    "fastgo/internal/service/user"
    "github.com/eininst/ninja"
    "github.com/gofiber/fiber/v2"
)

func init() {
    ninja.Provide(new(HellowordApi))
}

type HellowordApi struct {
    UserService user.UserService `inject:""`
}

func (h *HellowordApi) Add(c *fiber.Ctx) error {
    return c.SendString("hello word!")
}
```

> Inject resources...
```go
func Provide() {
    //inject resources
    rcli := data.NewRedisClient()
    ninja.Provide(rcli)
    ninja.Provide(rlock.New(rcli))
    ninja.Provide(data.NewRsClient(rcli))

    gormDB := data.NewDB()
    ninja.Provide(gormDB)

    //inject services
    ninja.Provide(user.NewUserService())
}
```

> Usage resources...
```go
package user

import (
    "fastgo/internal/code"
    "fastgo/internal/common/serr"
    "fastgo/internal/model"
    "fmt"
    "github.com/eininst/rlock"
    "github.com/go-redis/redis/v8"
    "github.com/jinzhu/copier"
    "gorm.io/gorm"
)

type UserService interface {
    Get(id int64) (*UserDTO, error)
}

type userService struct {
    DB          *gorm.DB      `inject:""`
    RedisClient *redis.Client `inject:""`
    Rlock       *rlock.Rlock  `inject:""`
}

func NewUserService() UserService {
    return &userService{}
}

func (us *userService) Get(id int64) (*UserDTO, error) {
    var u model.User
    us.DB.First(&u, id)
    if u.Id == 0 {
        msg := fmt.Sprintf("user is not found by %v", id)
        return nil, serr.NewServiceError(msg, code.ERROR_DATA_NOT_FOUND)
    }

    var udto UserDTO
    err := copier.Copy(&udto, &u)
    if err != nil {
        return nil, err
    }
    return &udto, nil
}

```


> See [examples](https://github.com/eininst/fastgo)

## License

*MIT*
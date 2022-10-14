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

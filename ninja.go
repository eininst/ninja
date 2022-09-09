package ninja

import (
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type App struct {
	*fiber.App
}

type Api interface {
	Router(r fiber.Router)
}

func New(config ...fiber.Config) *App {
	return &App{
		App: fiber.New(config...),
	}
}

func (a *App) Install(api Api) {
	Provide(api)
	Populate()
	api.Router(a.App)
}

func (a *App) Listen(addr string, config ...grace.Config) {
	if !strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf(":%s", addr)
	}
	grace.Listen(a.App, addr, config...)
}

func (a *App) ListenTLS(addr, certFile, keyFile string, config ...grace.Config) {
	if !strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf(":%s", addr)
	}
	grace.ListenTLS(a.App, addr, certFile, keyFile, config...)
}

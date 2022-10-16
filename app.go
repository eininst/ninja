package ninja

import (
	burst "github.com/eininst/fiber-middleware-burst"
	recovers "github.com/eininst/fiber-middleware-recover"
	redoc "github.com/eininst/fiber-middleware-redoc"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/ninja/serr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var (
	middlewareCfg = "middleware"
)

type Ninja struct {
	FiberApp *fiber.App
	Config   *Config
}

type Router interface {
	Init(router fiber.Router)
}

func New(path string, profile ...string) *Ninja {
	cfg := NewConfig(path, profile...)
	appConfig := cfg.GetAppConfig()
	app := fiber.New(fiber.Config{
		Prefork:      appConfig.Prefork,
		ReadTimeout:  time.Second * time.Duration(appConfig.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(appConfig.WriteTimeout),
		ErrorHandler: serr.ErrorHandler,
	})

	Provide(cfg, app)
	nj := &Ninja{FiberApp: app, Config: cfg}
	nj.useMiddleware()

	return nj
}

func (n *Ninja) Install(i Init) {
	Populate(i)
	i.Init()
}

func (n *Ninja) Router(r Router) {
	Populate(r)
	r.Init(n.FiberApp)
}

func (n *Ninja) App() *fiber.App {
	return n.FiberApp
}

func (n *Ninja) Listen(config ...grace.Config) {
	Populate()

	grace.Listen(n.FiberApp, n.Config.AppConfig.Port, config...)
}

func (n *Ninja) useMiddleware() {
	res := n.Config.Get("app", middlewareCfg)
	if res.IsArray() {
		midHandlerx := buildMiddleware()
		for _, r := range res.Array() {
			for k, v := range r.Map() {
				if handler, ok := midHandlerx[k]; ok {
					if v.Exists() {
						handler(n.App(), v)
					}
				}
			}
		}
	}
}

func buildMiddleware() map[string]func(app *fiber.App, value gjson.Result) {
	handlerx := map[string]func(app *fiber.App, value gjson.Result){}

	handlerx["recover"] = func(app *fiber.App, value gjson.Result) {
		rdefaltCfg := recovers.Config{}
		stackBuflenRes := value.Get("stackBuflen")
		if stackBuflenRes.Exists() {
			rdefaltCfg.StackTraceBufLen = int(stackBuflenRes.Int())
		}
		app.Use(recovers.New(rdefaltCfg))
	}
	handlerx["limiter"] = func(app *fiber.App, value gjson.Result) {
		rt := value.Get("rate").Int()
		bst := value.Get("burst").Int()
		timeout := value.Get("timeout").Int()
		if rt != 0 && bst != 0 && timeout != 0 {
			app.Use(burst.New(burst.Config{
				Limiter: rate.NewLimiter(rate.Limit(rt), int(bst)),
				Timeout: time.Second * time.Duration(timeout),
			}))
		}
	}

	handlerx["status"] = func(app *fiber.App, value gjson.Result) {
		path := "/status"
		pathRes := value.Get("path")
		if pathRes.String() != "" {
			path = pathRes.String()
		}
		app.Get(path, func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})
	}

	handlerx["logger"] = func(app *fiber.App, value gjson.Result) {
		f := "[Fiber] [${pid}] ${time} |${black}${status}|${latency}|${blue}${method} ${url}\n"
		tf := "2006/01/02 - 15:04:05"
		if value.Get("format").String() != "" {
			f = value.Get("format").String()
		}
		if value.Get("timeFormat").String() != "" {
			f = value.Get("timeFormat").String()
		}

		app.Use(logger.New(logger.Config{
			Format:     f,
			TimeFormat: tf,
		}))
	}

	handlerx["monitor"] = func(app *fiber.App, value gjson.Result) {
		app.Get("/metrics", monitor.New())
	}

	handlerx["swagger"] = func(app *fiber.App, value gjson.Result) {
		path := value.Get("path").String()
		json := value.Get("json").String()
		if path != "" && json != "" {
			app.Get(path, redoc.New(json))
		}
	}

	return handlerx
}

package http_server

import (
	"fmt"
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/kgantsov/limiter/pkg/limiter"
)

type App struct {
	api    huma.API
	router *fiber.App
	h      *Handler
	addr   int

	RateLimiter      *limiter.RateLimiter
	PathMap          map[string]string
	EnablePrometheus bool
}

func NewApp(
	addr int,
	RateLimiter *limiter.RateLimiter,
	PathMap map[string]string,
	EnablePrometheus bool,
) *App {

	router := fiber.New()
	api := humafiber.New(
		router, huma.DefaultConfig("Rate Limiter servie", "1.0.0"),
	)

	h := &Handler{
		RateLimiter: RateLimiter,
		PathMap:     PathMap,
	}

	h.ConfigureMiddleware(router)
	h.RegisterRoutes(api)

	app := &App{
		addr:             addr,
		api:              api,
		router:           router,
		h:                h,
		RateLimiter:      RateLimiter,
		PathMap:          PathMap,
		EnablePrometheus: EnablePrometheus,
	}

	return app
}

func (h *Handler) ConfigureMiddleware(router *fiber.App) {
	router.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05.999Z0700",
		TimeZone:   "Local",
		Format:     "${time} [INFO] ${locals:requestid} ${method} ${path} ${status} ${latency} ${error}​\n",
	}))

	router.Use(healthcheck.New())
	router.Use(helmet.New())

	router.Use(requestid.New())

	prometheus := fiberprometheus.New("limiter")
	prometheus.RegisterAt(router, "/metrics")
	router.Use(prometheus.Middleware)

	router.Get("/service/metrics", monitor.New())
	router.Use(recover.New())
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "rate-limiter-reduce",
			Method:      http.MethodPost,
			Path:        "/rate-limiters/:key/:max_tokens/:refill_time/:refill_amount/:tokens/",
			Summary:     "Reduce and get tokens",
			Description: "Reduce the number of tokens in the bucket and get the number of tokens left",
			Tags:        []string{"Rate Limiter"},
		},
		h.Reduce,
	)
	huma.Register(
		api,
		huma.Operation{
			OperationID: "rate-limiter-delete",
			Method:      http.MethodDelete,
			Path:        "/API/v1/rate-limiters/:key",
			Summary:     "Delete the key",
			Description: "Delete the key from the rate limiter",
			Tags:        []string{"Rate Limiter"},
		},
		h.Remove,
	)
	huma.Register(
		api,
		huma.Operation{
			OperationID: "rate-limiter-get-stats",
			Method:      http.MethodGet,
			Path:        "/API/v1/rate-limiters/stats",
			Summary:     "Stats of the rate limiter",
			Description: "Return the stats of the rate limiter service",
			Tags:        []string{"Rate Limiter"},
		},
		h.Stats,
	)
}

// Start starts the service.
func (a *App) Start() error {
	return a.router.Listen(fmt.Sprintf(":%d", a.addr))
}

// Close closes the service.
func (a *App) Close() {
	// s.e.Shutdown()
}

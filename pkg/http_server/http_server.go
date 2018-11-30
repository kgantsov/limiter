package http_server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/limiter/pkg/limiter"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	log "github.com/sirupsen/logrus"
)

type App struct {
	RateLimiter      *limiter.RateLimiter
	PathMap          map[string]string
	EnablePrometheus bool
}

type RateLimiterParams struct {
	MaxTokens    int64 `json:"max_tokens"`
	RefillTime   int64 `json:"refill_time"`
	RefillAmount int64 `json:"refill_amount"`
}

func DefineRoutes(app *App, r *gin.Engine) {
	r.GET("/stats/", app.Stat)
	v1 := r.Group("/API/v1")
	{
		v1.GET("/limiter/:key/:max_tokens/:refill_time/:refill_amount/:tokens/", app.ReduceLimiter)
	}

	if app.EnablePrometheus {
		for _, ri := range r.Routes() {
			app.PathMap[ri.Handler] = ri.Path
		}
	}
}

func ListenAndServe(app *App, port int, debug bool) {
	log.Infof("Strarting API service on a port: %d", port)

	if debug {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
	}

	r := gin.New()

	if app.EnablePrometheus {
		p := ginprometheus.NewPrometheus("gin")
		p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
			if path, ok := app.PathMap[c.HandlerName()]; ok {
				return path
			}

			return ""
		}

		p.Use(r)
	}

	r.Use(gin.Recovery())

	if debug {
		r.Use(gin.Logger())
	}

	DefineRoutes(app, r)

	r.Run(fmt.Sprintf(":%d", port))
}

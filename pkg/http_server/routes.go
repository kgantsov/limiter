package http_server

import "github.com/gin-gonic/gin"

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

package main

import (
	"flag"

	http_server "github.com/kgantsov/limiter/pkg/http_server"
	limiter "github.com/kgantsov/limiter/pkg/limiter"
	redis_server "github.com/kgantsov/limiter/pkg/redis_server"
)

func main() {
	httpPort := flag.Int("http_port", 9000, "HTTP Port")
	redisPort := flag.Int("redis_port", 46379, "Redis Port")
	// debug := flag.Bool("debug", false, "Debug flag")
	prometheus := flag.Bool("prometheus", false, "Enable prometheus")

	flag.Parse()

	rateLimiter := limiter.NewRateLimiter()
	go rateLimiter.StartCleanUpFullBuckets()

	app := http_server.NewApp(
		*httpPort, rateLimiter, make(map[string]string), *prometheus,
	)

	go redis_server.ListenAndServe(*redisPort, app.RateLimiter, *prometheus)
	// http_server.ListenAndServe(app, *httpPort, *debug)
	app.Start()
}

package main

import (
	"flag"

	http_server "github.com/kgantsov/limiter/pkg/http_server"
	limiter "github.com/kgantsov/limiter/pkg/limiter"
)

func main() {
	port := flag.Int("port", 9000, "Port")
	debug := flag.Bool("debug", false, "Debug flag")

	flag.Parse()

	app := &http_server.App{
		RateLimiter: limiter.NewRateLimiter(),
	}

	http_server.ListenAndServe(app, *port, *debug)
}

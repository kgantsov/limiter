package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	http_server "github.com/kgantsov/limiter/pkg/http_server"
	limiter "github.com/kgantsov/limiter/pkg/limiter"
	redis_server "github.com/kgantsov/limiter/pkg/redis_server"
)

func main() {
	httpPort := flag.Int("http_port", 9000, "HTTP Port")
	redisPort := flag.Int("redis_port", 46379, "Redis Port")
	debug := flag.Bool("debug", false, "Debug flag")

	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Default level for this example is info, unless debug flag is present
	log.Info().Msgf("Starting the application port: %d debug: %v", *httpPort, *debug)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	rateLimiter := limiter.NewRateLimiter()
	go rateLimiter.StartCleanUpFullBuckets()

	app := http_server.NewApp(
		*httpPort, rateLimiter, make(map[string]string),
	)

	server := redis_server.NewServer(*redisPort, app.RateLimiter)
	go server.ListenAndServe()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Info().Msgf("Received a signal: %d", sig)

		log.Info().Msg("Stopping the application")

		server.Stop()

		os.Exit(0)
	}()

	app.Start()
}

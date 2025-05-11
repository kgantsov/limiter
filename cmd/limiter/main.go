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
	logger "github.com/kgantsov/limiter/pkg/logger"
	redis_server "github.com/kgantsov/limiter/pkg/redis_server"
)

var (
	httpPort  int    // HTTP Port
	redisPort int    // HTTP Port
	logLevel  string // Log level
	logMode   string // Log mode

	cleanUpBucketsInterval = 300 // Clean up buckets interval in seconds
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	flag.IntVar(&httpPort, "http_port", 9000, "HTTP Port")
	flag.IntVar(&redisPort, "redis_port", 46379, "Redis Port")
	flag.StringVar(&logMode, "log_mode", "console", "Log mode: console, stackdriver")
	flag.StringVar(&logLevel, "log_level", "info", "Log level")
	flag.IntVar(&cleanUpBucketsInterval, "cleanup_interval", 300, "Clean up buckets interval in seconds")

	flag.Parse()

	logger.ConfigureLogger(logMode, logLevel)

	// Default level for this example is info, unless debug flag is present
	log.Info().Msgf("Starting the application port: %d", httpPort)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	rateLimiter := limiter.NewRateLimiter(time.Duration(cleanUpBucketsInterval) * time.Second)
	go rateLimiter.StartCleanUpFullBuckets()

	app := http_server.NewApp(
		httpPort, rateLimiter, make(map[string]string),
	)

	server := redis_server.NewServer(redisPort, app.RateLimiter)
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

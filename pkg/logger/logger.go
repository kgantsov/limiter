package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger(mode string, level string) {

	zerolog.TimeFieldFormat = time.RFC3339Nano

	if strings.ToUpper(mode) == "STACKDRIVER" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

		zerolog.LevelFieldName = "severity"
		zerolog.TimestampFieldName = "time"

		zerolog.LevelFieldMarshalFunc = func(level zerolog.Level) string {
			severity := map[zerolog.Level]string{
				zerolog.DebugLevel: "DEBUG",
				zerolog.InfoLevel:  "INFO",
				zerolog.WarnLevel:  "WARNING",
				zerolog.ErrorLevel: "ERROR",
				zerolog.FatalLevel: "CRITICAL",
				zerolog.PanicLevel: "EMERGENCY",
			}[level]
			return severity
		}

	} else {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339Nano,
			},
		)
	}

	logLevel, err := zerolog.ParseLevel(strings.ToUpper(level))
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}
}

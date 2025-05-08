//go:build !race

package logger

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigureLogger(t *testing.T) {
	ConfigureLogger("STACKDRIVER", "warning")

	assert.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
	assert.Equal(t, time.RFC3339Nano, zerolog.TimeFieldFormat)
	assert.Equal(t, "severity", zerolog.LevelFieldName)
	assert.Equal(t, "time", zerolog.TimestampFieldName)

	assert.Equal(t, "DEBUG", zerolog.LevelFieldMarshalFunc(zerolog.DebugLevel))
	assert.Equal(t, "INFO", zerolog.LevelFieldMarshalFunc(zerolog.InfoLevel))
	assert.Equal(t, "WARNING", zerolog.LevelFieldMarshalFunc(zerolog.WarnLevel))
	assert.Equal(t, "ERROR", zerolog.LevelFieldMarshalFunc(zerolog.ErrorLevel))
	assert.Equal(t, "CRITICAL", zerolog.LevelFieldMarshalFunc(zerolog.FatalLevel))
	assert.Equal(t, "EMERGENCY", zerolog.LevelFieldMarshalFunc(zerolog.PanicLevel))
}

package logging

import (
	"io"

	"{{CHIFRA}}/pkg/logger"
)

func Silence() func() {
	original := logger.GetLoggerWriter()
	logger.SetLoggerWriter(io.Discard)
	return func() {
		logger.SetLoggerWriter(original) // Restore original state
	}
}

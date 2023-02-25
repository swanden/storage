package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
)

const (
	twoFields = 2
)

// Compatibility check.
var _ logging.Logger = &Logger{lg: zerolog.Logger{}}

type diffStdWriter struct {
	out io.Writer
	err io.Writer
}

func (lw diffStdWriter) Write(p []byte) (int, error) {
	return lw.out.Write(p)
}

func (lw diffStdWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= zerolog.ErrorLevel {
		return lw.err.Write(p)
	}

	return lw.out.Write(p)
}

func New(ctx context.Context, options ...Option) (*Logger, error) {
	opts := GetDefaultOptions(ctx)

	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
	}

	if err := Validate(opts); err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logg := (zerolog.New(diffStdWriter{out: os.Stdout, err: os.Stderr})).
		Level(opts.rawLevel).
		Sample(&zerolog.BasicSampler{N: uint32(opts.basicSamplerNum)}).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{lg: logg}, nil
}

type Logger struct {
	lg zerolog.Logger
}

func (logger *Logger) Level(levelRaw string) *Logger {
	level, _ := zerolog.ParseLevel(levelRaw)
	lg := logger.lg.Level(level)
	logger.lg = lg

	return logger
}

func (logger *Logger) Debug() Event {
	return Event{Event: logger.lg.Debug()}
}

func (logger *Logger) Info() Event {
	return Event{Event: logger.lg.Info()}
}

func (logger *Logger) Error() Event {
	return Event{Event: logger.lg.Error()}
}

func (logger *Logger) Panic() Event {
	return Event{Event: logger.lg.Panic()}
}

func (logger *Logger) Fatal() Event {
	return Event{Event: logger.lg.Fatal()}
}

// Log implements the logging.Logger interface.
// got from github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 .
func (logger *Logger) Log(lvl logging.Level, msg string) {
	switch lvl {
	case logging.DEBUG:
		logger.lg.Debug().Msg(msg)
	case logging.INFO:
		logger.lg.Info().Msg(msg)
	case logging.WARNING:
		logger.lg.Warn().Msg(msg)
	case logging.ERROR:
		logger.lg.Error().Msg(msg)
	default:
	}
}

func (logger *Logger) Printf(a string, b ...interface{}) {
	logger.lg.Info().Msg(fmt.Sprintf(a, b...))
}

func (logger *Logger) Debugf(a string, b ...interface{}) {
	logger.lg.Debug().Msg(fmt.Sprintf(a, b...))
}

func (logger *Logger) Infof(a string, b ...interface{}) {
	logger.lg.Info().Msg(fmt.Sprintf(a, b...))
}

func (logger *Logger) Errorf(a string, b ...interface{}) {
	logger.lg.Error().Msg(fmt.Sprintf(a, b...))
}

// With implements the logging.Logger interface.
// got from github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 .
func (logger *Logger) With(fields ...string) logging.Logger {
	vals := make(map[string]interface{}, len(fields)/twoFields)
	for i := 0; i < len(fields); i += 2 {
		vals[fields[i]] = fields[i+1]
	}

	return &Logger{lg: logger.lg.With().Fields(vals).Logger()}
}

type Event struct {
	*zerolog.Event
}

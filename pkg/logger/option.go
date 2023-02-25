package logger

import (
	"log"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

const (
	samplerNumDef = 1
)

type Option func(*Options) error

type Options struct {
	rawLevel        zerolog.Level
	basicSamplerNum int64
}

func GetDefaultOptions(_ context.Context) Options {
	return Options{
		rawLevel:        zerolog.InfoLevel,
		basicSamplerNum: samplerNumDef,
	}
}

func Validate(opts Options) error {
	if opts.basicSamplerNum == 0 {
		return ErrSamplerNum
	}

	return nil
}

func WithLevel(rawLevel string) Option {
	return func(options *Options) error {
		level, err := zerolog.ParseLevel(rawLevel)
		if err != nil {
			log.Panicf("Bad log level: %v, error: %v\n", rawLevel, err)
		}

		options.rawLevel = level

		return nil
	}
}

func WithBasicSamplerNum(num int64) Option {
	return func(options *Options) error {
		options.basicSamplerNum = num

		return nil
	}
}

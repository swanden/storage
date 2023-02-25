package conf

import (
	"log"
	"os"
	"strconv"
	"time"
)

func hasValue(name string) (string, bool) {
	value := os.Getenv(name) // Package conf is a configuration tool
	if value == "" {
		return "", false
	}

	return value, true
}

func IntValue(name string, defaultValue int) int {
	if value, has := hasValue(name); has {
		v, err := strconv.Atoi(value)
		if err != nil {
			panicBadEnvKey(name, err)
		}

		return v
	}

	return defaultValue
}

func StrValue(name, defaultValue string) string {
	if value, has := hasValue(name); has {
		return value
	}

	return defaultValue
}

func StrValueRequired(name string) string {
	if value, has := hasValue(name); has {
		return value
	}

	panicMissReqKey(name)

	return ""
}

func TimeDurValue(name string, defaultValue time.Duration) time.Duration {
	if value, has := hasValue(name); has {
		v, err := strconv.Atoi(value)
		if err != nil {
			panicBadEnvKey(name, err)
		}

		return time.Duration(v) * time.Millisecond
	}

	return defaultValue
}

func panicBadEnvKey(name string, err error) {
	log.Panicf("Bad env key: %v, error: %v", name, err)
}

func panicMissReqKey(name string) {
	log.Panicf("Missed required env key: %v", name)
}

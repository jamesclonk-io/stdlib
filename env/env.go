package env

import (
	"os"
)

func Get(key string, nvl string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nvl
	}
	return value
}

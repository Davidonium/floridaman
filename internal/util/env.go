package util

import "os"

// GetEnvDefault gets the `key` environment variable or returns the default value.
func GetEnvDefault(key, d string) string {
	e, ok := os.LookupEnv(key)

	if !ok {
		e = d
	}

	return e
}

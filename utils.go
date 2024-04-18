package tracer

import (
	"os"
	"strconv"
)

// GetterBool -
func GetterBool(key string, defaultValue bool) bool {
	env, ok := os.LookupEnv(key)
	if ok {
		res, err := strconv.ParseBool(env)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

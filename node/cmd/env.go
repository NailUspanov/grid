package main

import (
	"os"
	"strconv"
)

var (
	httpPort = GetterInt("HTTP_PORT", 8000)
)

func Getter(key, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	return defaultValue
}

func GetterInt(key string, defaultValue int) int {
	env, ok := os.LookupEnv(key)
	if ok {
		res, err := strconv.ParseInt(env, 10, 32)
		if err == nil {
			return int(res)
		}
	}
	return defaultValue
}

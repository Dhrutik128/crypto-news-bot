package config


import (
	"os"
	"strconv"
)

func String(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Int(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}

func Float64(key string, fallback float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}

func Bool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}

func Float32(key string, fallback float64) float32 {
	return float32(Float64(key, fallback))
}


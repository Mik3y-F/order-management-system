package pkg

import "os"

func MustGetEnv(key string) string {
	if v := os.Getenv(key); v == "" {
		panic("Missing required environment variable: " + key)
	} else {
		return v
	}
}

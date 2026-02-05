package config

import "os"

type Config struct {
	Addr        string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return Config{
		Addr:        addr,
		DatabaseURL: must("DATABASE_URL"),
		JWTSecret:   must("JWT_SECRET"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}

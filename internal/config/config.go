package config

import "os"

type Config struct {
	Port   string
	DBPath string
}

func Load() Config {
	c := Config{
		Port:   ":8080",
		DBPath: "flaggy.db",
	}
	if v := os.Getenv("FLAGGY_PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("FLAGGY_DB_PATH"); v != "" {
		c.DBPath = v
	}
	return c
}

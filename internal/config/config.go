package config

import "os"

type Config struct {
	Port      string
	DBPath    string
	MasterKey string // Required for admin routes (key management, flag CRUD)
	CORSEnabled bool
}

func Load() Config {
	c := Config{
		Port:        ":8080",
		DBPath:      "flaggy.db",
		MasterKey:   os.Getenv("FLAGGY_MASTER_KEY"),
		CORSEnabled: os.Getenv("FLAGGY_CORS") != "false",
	}
	if v := os.Getenv("FLAGGY_PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("FLAGGY_DB_PATH"); v != "" {
		c.DBPath = v
	}
	return c
}

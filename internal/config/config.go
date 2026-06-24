package config

import (
"fmt"
"os"

"github.com/joho/godotenv"
)

type Config struct {
AppPort string
AppEnv  string
}

func Load() (*Config, error) {
_ = godotenv.Load()

cfg := &Config{
AppPort: getEnv("APP_PORT", "8080"),
AppEnv:  getEnv("APP_ENV", "local"),
}

if cfg.AppPort == "" {
return nil, fmt.Errorf("APP_PORT is required")
}

return cfg, nil
}

func getEnv(key string, defaultValue string) string {
value := os.Getenv(key)
if value == "" {
return defaultValue
}

return value
}

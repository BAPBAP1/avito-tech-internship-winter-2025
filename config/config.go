
package config

import "os"

type Config struct {
 DBHost     string
 DBPort     string
 DBUser     string
 DBPassword string
 DBName     string
 JWTSecret  string
}
func Load() *Config {
	return &Config{
	 DBHost:     getEnv("DB_HOST", "localhost"),
	 DBPort:     getEnv("DB_PORT", "5432"),
	 DBUser:     getEnv("DB_USER", "postgres"),
	 DBPassword: getEnv("DB_PASSWORD", "postgres"),
	 DBName:     getEnv("DB_NAME", "avito_merch"),
	 JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}
   }
   
   func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
	 return value
	}
	return defaultVal
   }
   
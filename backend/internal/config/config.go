package config

import "os"

type KeycloakConfig struct {
	AuthServerURL string // базовый URL Keycloak, без /realms (например: https://keycloak.example.com/auth)
	Realm         string
	ClientID      string
	ClientSecret  string
}

type Config struct {
	PostgresDSN string
	MongoURI    string
	Keycloak    KeycloakConfig
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		PostgresDSN: getenv("APP_POSTGRES_DSN", ""),
		MongoURI:    getenv("APP_MONGO_URI", ""),
		Keycloak: KeycloakConfig{
			AuthServerURL: getenv("APP_KEYCLOAK_URL", ""), // напр. https://keycloak.example.com/auth
			Realm:         getenv("APP_KEYCLOAK_REALM", ""),
			ClientID:      getenv("APP_KEYCLOAK_CLIENT_ID", ""),
			ClientSecret:  getenv("APP_KEYCLOAK_CLIENT_SECRET", ""),
		},
	}
}

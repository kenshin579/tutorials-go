package config

type Config struct {
	Server   ServerConfig
	Keycloak KeycloakConfig
}

type ServerConfig struct {
	Port string
}

type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8081",
		},
		Keycloak: KeycloakConfig{
			BaseURL:      "http://localhost:8080",
			Realm:        "myrealm",
			ClientID:     "mybackend",
			ClientSecret: "your-client-secret", // TODO: Replace with actual client secret
		},
	}
}

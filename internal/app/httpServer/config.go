package httpServer

type Config struct {
	BindAddr    string
	DatabaseURL string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		DatabaseURL: "././identifier.sqlite",
	}
}

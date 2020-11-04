package telegrambot

// Config ...
type Config struct {
	Title    string
	Telegram telegram           `toml:"telegram"`
	Serial   serialcomunication `toml:"serialcomunication"`
	LogLevel string             `toml:"log_level"`
}

type telegram struct {
	APIURL string `toml:"api_url"`
}

type serialcomunication struct {
	SerialName  string `toml:"serial_name"`
	SerialSpeed int    `toml:"serial_speed"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
	}
}

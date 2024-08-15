package config

type Config struct {
	LogLevel *string
}

func NewConfig(flags *Flags) *Config {
	return &Config{
		LogLevel: flags.LogLevel,
	}
}

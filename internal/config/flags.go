package config

const (
	DefaultLogLevel = "info"
)

type Flags struct {
	LogLevel *string
}

func NewFlags() *Flags {
	return &Flags{
		LogLevel: strPtr(DefaultLogLevel),
	}
}

func strPtr(s string) *string {
	return &s
}

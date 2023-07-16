package auth

import (
	"github.com/spf13/pflag"

	"github.com/triabokon/gotagv/internal/flags"
)

type Config struct {
	JWTSecret string
}

func (c *Config) Flags(prefix string) *pflag.FlagSet {
	const name = "AuthConfig"
	f := pflag.NewFlagSet(name, pflag.PanicOnError)

	f.StringVar(&c.JWTSecret, "jwt_secret", "", "secret to generate jwt tokens")
	return flags.MapWithPrefix(f, name, pflag.PanicOnError, prefix)
}

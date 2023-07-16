package server

import (
	"github.com/spf13/pflag"

	"github.com/triabokon/gotagv/internal/auth"
	"github.com/triabokon/gotagv/internal/postgresql"
	"github.com/triabokon/gotagv/internal/server"
)

type Config struct {
	HTTP     server.Config
	Postgres postgresql.Config

	Auth auth.Config
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("GoTagVConfig", pflag.PanicOnError)

	f.AddFlagSet(c.HTTP.Flags("http"))
	f.AddFlagSet(c.Postgres.Flags("postgres"))

	f.AddFlagSet(c.Auth.Flags("auth"))
	return f
}

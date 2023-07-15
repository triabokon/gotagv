package server

import (
	"github.com/spf13/pflag"

	"github.com/triabokon/gotagv/server"
)

type Config struct {
	HTTP server.Config
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("GoTagVConfig", pflag.PanicOnError)

	f.AddFlagSet(c.HTTP.Flags("http"))

	return f
}

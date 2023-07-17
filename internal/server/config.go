package server

import (
	"time"

	"github.com/spf13/pflag"

	"github.com/triabokon/gotagv/internal/flags"
)

type Config struct {
	Bind         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c *Config) Flags(prefix string) *pflag.FlagSet {
	const name = "HTTPConfig"
	f := pflag.NewFlagSet(name, pflag.PanicOnError)

	f.StringVar(&c.Bind, "bind", "0.0.0.0:8080", "bind address for http server")
	f.DurationVar(&c.ReadTimeout, "read_timeout", 5*time.Second, "read timeout as described in net/http.Server")
	f.DurationVar(&c.WriteTimeout, "write_timeout", 5*time.Second, "write timeout as described in net/http.Server")
	f.DurationVar(&c.IdleTimeout, "idle_timeout", time.Minute, "idle timeout as described in net/http.Server")
	return flags.MapWithPrefix(f, name, pflag.PanicOnError, prefix)
}

package server

import (
	"github.com/MyChaOS87/patAi/config"
	"github.com/MyChaOS87/patAi/internal/api/router"
)

type (
	Option func(*Config)
	Config struct {
		api          *config.APIConfig
		childRouters []router.Router
	}
)

func newDefaultConfig() *Config {
	return &Config{}
}

func API(apiConfig *config.APIConfig) Option {
	return func(c *Config) {
		c.api = apiConfig
	}
}

func ChildRouters(routers ...router.Router) Option {
	return func(c *Config) {
		c.childRouters = append(c.childRouters, routers...)
	}
}

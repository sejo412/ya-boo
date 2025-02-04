package app

import "github.com/sejo412/ya-boo/pkg/config"

type app struct {
	cfg *config.Config
}

func newApp(cfg *config.Config) *app {
	return &app{cfg}
}

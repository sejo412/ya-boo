package app

import (
	"log"

	"github.com/sejo412/ya-boo/pkg/config"
)

type App struct {
	Cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg}
}

func (a *App) Run() error {
	log.Println("starting server")

	log.Printf("config: %#v\n", a.Cfg)
	return nil
}

package app

import (
	"context"
	"log"

	"github.com/sejo412/ya-boo/pkg/config"
	"github.com/sejo412/ya-boo/pkg/tg"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tgBot, err := tg.NewBot(a.Cfg.TgSecret)
	if err != nil {
		return err
	}

	tgBot.Bot.Start(ctx)

	return nil
}

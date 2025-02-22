package app

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/sejo412/ya-boo/pkg/ai"
	"github.com/sejo412/ya-boo/pkg/config"
)

type App struct {
	cfg       *config.Config
	telegram  *bot.Bot
	aiClients map[string]*ai.Client
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg, telegram: &bot.Bot{}, aiClients: make(map[string]*ai.Client)}
}

func (a *App) Run() error {
	log.Println("starting server")

	log.Printf("config: %#v\n", a.cfg)

	ctxTg, cancelTg := context.WithCancel(context.Background())
	defer cancelTg()

	localAi := ai.NewClient("http://127.0.0.1:8000", "")
	a.aiClients["local"] = localAi

	return a.StartTelegramBot(ctxTg, a.cfg.TgSecret)
}

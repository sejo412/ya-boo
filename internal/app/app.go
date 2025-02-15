package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sejo412/ya-boo/pkg/ai"
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

	ctxTg, cancelTg := context.WithCancel(context.Background())
	defer cancelTg()

	tgBot, err := tg.NewBot(a.Cfg.TgSecret)
	if err != nil {
		return err
	}

	go tgBot.Start(ctxTg)

	fmt.Println("Bot started")

	for ch := range len(tg.Updates) {
		fmt.Println("Updates:", ch)
	}

	time.Sleep(60 * time.Minute)

	ctxAi, cancelAi := context.WithCancel(context.Background())
	defer cancelAi()

	aiLocal := ai.NewAI("http://127.0.0.1:8000", "")
	fwCh := make(chan string)
	startAi(ctxAi, aiLocal, fwCh)

	return nil
}

func startAi(ctx context.Context, llm *ai.AI, ch <-chan string) {

}

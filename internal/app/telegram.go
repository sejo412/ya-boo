package app

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) StartTelegramBot(ctx context.Context, token string) error {
	var err error
	opts := []bot.Option{
		bot.WithDefaultHandler(a.handler),
		bot.WithDebug(),
	}
	a.telegram, err = bot.New(token, opts...)
	if err != nil {
		return fmt.Errorf("error init bot: %w", err)
	}
	a.telegram.Start(ctx)
	return nil
}

func (a *App) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		resp, err := a.aiClients["local"].ChatCompletion(ctx, update.Message.Text)
		if err != nil {
			log.Printf("error completing message: %v", err)
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf(resp),
		})
		if err != nil {
			log.Printf("error sending message: %v", err)
		}
	}
}

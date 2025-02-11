package tg

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	*bot.Bot
}

func NewBot(token string) (*Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithDebug(),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("error init bot: %w", err)
	}
	return &Bot{Bot: b}, nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("echo: %s", update.Message.Text),
		})
		if err != nil {
			log.Printf("error sending message: %v", err)
		}
	}
}

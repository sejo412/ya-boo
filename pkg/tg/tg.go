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
	defaultHandlerFunc HandlerFunc
}

type HandlerFunc func(ctx context.Context, bot *Bot, update *models.Update, ch chan<- struct{})

var Updates = make(chan *models.Update)

type Option func(*Bot)

func WithDefaultHandler(handler HandlerFunc) Option {
	return func(b *Bot) {
		b.defaultHandlerFunc = handler
	}
}

func NewBot(token string) (*Bot, error) {
	bot.WithDefaultHandler(handler)
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

/*
func (b *Bot) Start(ctx context.Context) {
	b.Bot.Start(ctx)
}

*/

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("echo: %s", update.Message.Text),
		})
		if err != nil {
			log.Printf("error sending message: %v", err)
		}
		Updates <- update

	}
}

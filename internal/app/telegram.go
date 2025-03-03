package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	m "github.com/sejo412/ya-boo/pkg/models"
)

func (a *App) StartTelegramBot(ctx context.Context) error {
	var err error
	initMode, err := a.isInitMode(ctx)
	if err != nil {
		return fmt.Errorf("error check init mode: %w", err)
	}
	opts := []bot.Option{
		bot.WithMiddlewares(a.checkUser),
		bot.WithDefaultHandler(a.defaultHandler),
		bot.WithDebug(),
	}
	a.telegram, err = bot.New(a.cfg.TgSecret, opts...)
	if err != nil {
		return fmt.Errorf("error init bot: %w", err)
	}
	if initMode {
		a.initID = a.telegram.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, a.initHandler)
	}
	a.telegram.RegisterHandler(bot.HandlerTypeMessageText, "/", bot.MatchTypePrefix, a.commandHandler)
	a.telegram.Start(ctx)
	return nil
}

func (a *App) initHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	initSuccess := false
	if update.Message != nil {
		var err error
		var resp string
		switch update.Message.Text {
		case fmt.Sprintf("/init %s", a.cfg.InitBotSecret):
			if err = cmdInitFirstAdmin(ctx, a.db, m.User{
				User: &models.User{
					ID:        update.Message.From.ID,
					FirstName: update.Message.From.FirstName,
					LastName:  update.Message.From.LastName,
					Username:  update.Message.From.Username,
				},
				Role: m.RoleAdmin,
			}); err != nil {
				resp = err.Error()
			} else {
				resp = MessageInitOk
				initSuccess = true
			}
		default:
			resp = MessageUnknownCommand + "or" + MessageBadInitSecret
			resp += fmt.Sprintf(MessageInit, update.Message.From.ID, update.Message.From.Username,
				update.Message.From.FirstName, update.Message.From.LastName)
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      bot.EscapeMarkdown(resp),
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			log.Printf("[initHandler] error sending message: %v", err)
		}
	}
	if initSuccess {
		log.Printf("[initHandler] init success")
		a.telegram.UnregisterHandler(a.initID)
		a.initID = ""
	}
}

func (a *App) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		var err error
		resp, err := a.aiClients["local"].ChatCompletion(ctx, update.Message.Text)
		if err != nil {
			log.Printf("error completing message: %v", err)
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      bot.EscapeMarkdown(resp),
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			log.Printf("[defaultHandler] error sending message: %v", err)
		}
	}
}

func (a *App) commandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		var err error
		var resp, res string
		var id int
		splited := strings.Split(update.Message.Text, " ")
		switch splited[0] {
		case CmdList.String():
			res, err = cmdListUsers(ctx, a.db)
			if err != nil {
				resp = err.Error()
				log.Printf("[commandHandler] error listing users: %v", err)
			} else {
				resp = res
			}
		case CmdApprove.String():
			id, err = strconv.Atoi(splited[1])
			if err != nil {
				log.Printf("[commandHandler] error approving user: %v", err)
			}
			res, err = cmdApproveUser(ctx, a.db, m.User{
				User: &models.User{
					ID: int64(id),
				},
				Role: m.RoleRegularUser,
			})
			if err != nil {
				resp = err.Error()
				log.Printf("[commandHandler] error approving user: %v", err)
			} else {
				resp = res
			}
		default:
			resp = MessageUnknownCommand
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      bot.EscapeMarkdown(resp),
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			log.Printf("[commandHandler] error sending message: %v", err)
		}
	}
}

func (a *App) isInitMode(ctx context.Context) (bool, error) {
	adminPresents, err := a.db.IsAdminsInitialized(ctx)
	if err != nil {
		return false, err
	}
	return !adminPresents, nil
}

func (a *App) checkUser(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if a.initID != "" {
			next(ctx, b, update)
			return
		}
		if update.Message != nil {
			var resp = "you not authorized send messages to bot"
			reg, err := a.db.IsRegisteredUser(ctx, update.Message.From.ID)
			if err != nil {
				log.Printf("[registeredUser] error check user: %v", err)
				resp = fmt.Sprintf("error check user: %v", err)
			}
			if reg {
				next(ctx, b, update)
				return
			}
			visited, err := a.db.IsUserPresent(ctx, update.Message.From.ID)
			if err != nil {
				log.Printf("[registeredUser] error check user: %v", err)
				resp = fmt.Sprintf("error check user: %v", err)
			}
			if !visited {
				err = a.db.UpsertUser(ctx, m.User{
					User: &models.User{
						ID:        update.Message.From.ID,
						FirstName: update.Message.From.FirstName,
						LastName:  update.Message.From.LastName,
						Username:  update.Message.From.Username,
					},
					Role: m.RoleUnknown,
				})
				if err != nil {
					log.Printf("[registeredUser] error upserting user: %v", err)
					resp = fmt.Sprintf("error upserting user: %v", err)
				}
			}
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      resp,
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				log.Printf("[checkUser] error sending message: %v", err)
			}
			return
		}
	}
}

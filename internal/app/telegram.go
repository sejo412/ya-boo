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
				LLM:  m.LLM{Id: 0},
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
		var resp string
		llm, err := a.db.GetUserLLM(ctx, update.Message.From.ID)
		if err != nil || llm.Id == 0 {
			resp = MessageHelper
		} else {
			resp, err = a.aiClients[llm.Id].ChatCompletion(ctx, update.Message.Text)
			if err != nil {
				log.Printf("error completing message: %v", err)
				resp = MessageLLMError
			}
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
		var isAdmin bool
		splited := strings.Split(update.Message.Text, " ")
		switch splited[0] {
		case CmdList.String():
			isAdmin, err = a.db.IsAdmin(ctx, update.Message.From.ID)
			if err != nil || !isAdmin {
				resp = MessageNotAuthorized
				break
			}
			res, err = cmdListUsers(ctx, a.db)
			if err != nil {
				resp = err.Error()
				log.Printf("[commandHandler] error listing users: %v", err)
			} else {
				resp = res
			}
		case CmdApprove.String():
			isAdmin, err = a.db.IsAdmin(ctx, update.Message.From.ID)
			if err != nil || !isAdmin {
				resp = MessageNotAuthorized
				break
			}
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
		case CmdBan.String():
			isAdmin, err = a.db.IsAdmin(ctx, update.Message.From.ID)
			if err != nil || !isAdmin {
				resp = MessageNotAuthorized
				break
			}
			id, err = strconv.Atoi(splited[1])
			if err != nil {
				log.Printf("[commandHandler] error ban user: %v", err)
			}
			res, err = cmdBanUser(ctx, a.db, m.User{
				User: &models.User{
					ID: int64(id),
				},
			})
			if err != nil {
				resp = err.Error()
				log.Printf("[commandHandler] error ban user: %v", err)
			} else {
				resp = res
			}
		case CmdLlmAdd.String():
			isAdmin, err = a.db.IsAdmin(ctx, update.Message.From.ID)
			if err != nil || !isAdmin {
				resp = MessageNotAuthorized
				break
			}
			if err = cmdLlmAdd(ctx, a.db, update.Message.Text); err != nil {
				resp = err.Error()
			}
			resp = MessageLLMAddSuccess
		case CmdLlmRemove.String():
			isAdmin, err = a.db.IsAdmin(ctx, update.Message.From.ID)
			if err != nil || !isAdmin {
				resp = MessageNotAuthorized
				break
			}
			id, err = strconv.Atoi(splited[1])
			if err != nil {
				log.Printf("[commandHandler] error remove llm: %v", err)
				resp = fmt.Sprintf("error remov llm: %v", err)
				break
			}
			if err = cmdLlmRemove(ctx, a.db, int64(id)); err != nil {
				log.Printf("[commandHandler] error remove llm: %v", err)
				resp = fmt.Sprintf("error remove llm: %v", err)
				break
			}
		case CmdLlmList.String():
			res, err = cmdLlmList(ctx, a.db)
			if err != nil {
				resp = err.Error()
				log.Printf("[commandHandler] error listing llm: %v", err)
			} else {
				resp = res
			}
		case CmdLlmUse.String():
			id, err = strconv.Atoi(splited[1])
			if err != nil {
				log.Printf("[commandHandler] error use llm: %v", err)
				resp = fmt.Sprintf("error use llm: %v", err)
				break
			}
			if err = cmdLlmUse(ctx, a.db, update.Message.From.ID, int64(id)); err != nil {
				resp = err.Error()
			}
			resp = fmt.Sprintf("switched to llm: %d", id)
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

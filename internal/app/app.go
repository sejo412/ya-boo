package app

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/sejo412/ya-boo/pkg/ai"
	"github.com/sejo412/ya-boo/pkg/config"
	"github.com/sejo412/ya-boo/pkg/models"
)

type App struct {
	cfg       *config.Config
	telegram  *bot.Bot
	aiClients map[int64]AiClient
	db        Storage
	initID    string
}

type AiClient interface {
	ChatCompletion(ctx context.Context, req string) (resp string, err error)
}

type Storage interface {
	Open(dsn string) error
	Close()
	Ping() error
	IsAdminsInitialized(ctx context.Context) (bool, error)
	IsUserPresent(ctx context.Context, id int64) (bool, error)
	IsAdmin(ctx context.Context, id int64) bool
	IsRegisteredUser(ctx context.Context, id int64) (bool, error)
	IsWaitingApprove(ctx context.Context, id int64) (bool, error)
	UpsertUser(ctx context.Context, user models.User) error
	UpdateUserRole(ctx context.Context, user models.User, role models.Role) error
	ListUsers(ctx context.Context) ([]models.User, error)
	GetLLMs(ctx context.Context) ([]models.LLM, error)
	GetUserLLM(ctx context.Context, userID int64) (models.LLM, error)
	AddLLM(ctx context.Context, llm models.LLM) error
	RemoveLLM(ctx context.Context, llm models.LLM) error
	SetUserLLM(ctx context.Context, userID int64, llm models.LLM) error
}

func NewApp(cfg *config.Config, storage Storage) *App {
	return &App{
		cfg:       cfg,
		telegram:  &bot.Bot{},
		aiClients: make(map[int64]AiClient),
		db:        storage,
		initID:    "",
	}
}

func (a *App) Run() error {
	log.Println("starting server")
	log.Printf("config: %#v\n", a.cfg)
	if err := a.db.Open(a.cfg.Dsn); err != nil {
		return err
	}
	defer a.db.Close()

	if err := a.initLLMs(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return a.StartTelegramBot(ctx)
}

func (a *App) initLLMs() error {
	llms, err := a.db.GetLLMs(context.Background())
	if err != nil {
		return fmt.Errorf("error get llms: %w", err)
	}
	for _, llm := range llms {
		a.aiClients[llm.ID] = ai.NewClient(llm.Endpoint, llm.Token)
	}
	return nil
}

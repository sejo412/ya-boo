package app

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/sejo412/ya-boo/pkg/ai"
	"github.com/sejo412/ya-boo/pkg/config"
	"github.com/sejo412/ya-boo/pkg/models"
)

type App struct {
	cfg       *config.Config
	telegram  *bot.Bot
	aiClients map[string]*ai.Client
	db        Storage
	initID    string
}

type Storage interface {
	Open(dsn string) error
	Close()
	Ping() error
	IsAdminsInitialized(ctx context.Context) (bool, error)
	IsUserPresent(ctx context.Context, id int64) (bool, error)
	IsRegisteredUser(ctx context.Context, id int64) (bool, error)
	IsWaitingApprove(ctx context.Context, id int64) (bool, error)
	UpsertUser(ctx context.Context, user models.User) error
	UpdateUserRole(ctx context.Context, user models.User, role models.Role) error
	ListUsers(ctx context.Context) ([]models.User, error)
}

func NewApp(cfg *config.Config, storage Storage) *App {
	return &App{
		cfg:       cfg,
		telegram:  &bot.Bot{},
		aiClients: make(map[string]*ai.Client),
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

	localAi := ai.NewClient("http://127.0.0.1:8000", "")
	a.aiClients["local"] = localAi

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return a.StartTelegramBot(ctx)
}

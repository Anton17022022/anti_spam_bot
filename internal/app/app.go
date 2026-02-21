package app

import (
	"context"
	"fmt"
	"log"

	antispambot "telegram-antispam-bot/internal/app/anti_spam_bot"
	"telegram-antispam-bot/internal/infra/config"
	"telegram-antispam-bot/internal/infra/storage"
	models_err_app "telegram-antispam-bot/internal/models/errors/app"
)

// App is app
type App struct {
	service service
	infra   infra
}

type infra struct {
	conf    *config.Config
	storage *storage.Storage
}

type service struct {
	antiSpamBot antispambot.Bot
}

// NewApp init app components. Return instance.
func NewApp(ctx context.Context) (*App, error) {
	a := App{}

	// not err return. init internal. not ideomatic
	if err := a.initInfra(); err != nil {
		return nil, fmt.Errorf("%w:%v", models_err_app.ErrInitApp, err)
	}

	if err := a.initService(ctx); err != nil {
		return nil, fmt.Errorf("%w:%v", models_err_app.ErrInitApp, err)
	}

	return &a, nil
}

func (a *App) initInfra() error {
	conf, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config.NewConfig: %w", err)
	}

	storage, err := storage.NewStorage(conf)
	if err != nil {
		return fmt.Errorf("storage.NewStorage: %w", err)
	}

	a.infra = infra{
		conf:    conf,
		storage: storage,
	}

	return nil
}

func (a *App) initService(ctx context.Context) error {
	bot, err := antispambot.NewAntiSpamBot(a.infra.conf, a.infra.storage)
	if err != nil {
		return fmt.Errorf("%w:%v", models_err_app.ErrInitService, err)
	}

	log.Printf("Authorized on account %s", bot.UserName)

	bot.RegisterRoutes(ctx)

	a.service = service{antiSpamBot: *bot}

	return nil
}

// ListenAndServe start app
func (a *App) ListenAndServe() {
	go func() {
		a.service.antiSpamBot.Start()
	}()

	a.service.antiSpamBot.StartDelSpamMessage()
}

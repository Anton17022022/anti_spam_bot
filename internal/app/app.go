package app

import (
	"fmt"
	"log"

	antispambot "telegram-antispam-bot/internal/app/anti_spam_bot"
	"telegram-antispam-bot/internal/infra/config"
	models_err_app "telegram-antispam-bot/internal/models/errors/app"
)

// App is app
type App struct {
	service service
	infra   infra
}

type infra struct {
	conf *config.Config
}

type service struct {
	antiSpamBot antispambot.Bot
}

// NewApp init app components. Return instance.
func NewApp() (*App, error) {
	a := App{}

	// not err return. init internal. not ideomatic
	if err := a.initInfra(); err != nil {
		return nil, fmt.Errorf("%w:%v", models_err_app.ErrInitApp, err)
	}

	if err := a.initService(); err != nil {
		return nil, fmt.Errorf("%w:%v", models_err_app.ErrInitApp, err)
	}

	return &a, nil
}

func (a *App) initInfra() error {
	conf, err := config.NewConfig()
	if err != nil {
		return err
	}

	a.infra = infra{conf: conf}

	return nil
}

func (a *App) initService() error {
	bot, err := antispambot.NewAntiSpamBot(a.infra.conf)
	if err != nil {
		return fmt.Errorf("%w:%v", models_err_app.ErrInitService, err)
	}

	log.Printf("Authorized on account %s", bot.UserName)

	a.service = service{antiSpamBot: *bot}

	return nil
}

// ListenAndServe start app
func (a *App) ListenAndServe() {
	a.service.antiSpamBot.StartDelSpamMessage()
}

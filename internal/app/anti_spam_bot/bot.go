package antispambot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"

	"telegram-antispam-bot/internal/infra/config"
	admmodels "telegram-antispam-bot/internal/models/adm_models"
	models_errors_anti_spambot "telegram-antispam-bot/internal/models/errors/anti_spam_bot"
)

// antiSpamBot is API to
type antiSpamBot interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

// Storage is interface for database
type Storage interface {
	DelWordFromBadWords(word string) error
	GetListBadWords() ([]string, error)
	InsertWordToBadWords(word string) error
}

// Bot is bot for struggle spam
type Bot struct {
	Bot      antiSpamBot
	BotAdm   *telebot.Bot
	UserName string
	Storage  Storage
	conf     *config.Config
}

// NewAntiSpamBot is constructor AntiSpamBot. Return new instance.
func NewAntiSpamBot(conf *config.Config, storage Storage) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(conf.BotAntiSpam.Settings.Token)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", models_errors_anti_spambot.ErrInitBot, err)
	}

	botAdm, err := telebot.NewBot(telebot.Settings{
		Token:  conf.BotAntiSpam.Settings.AdmToken,
		Poller: &telebot.LongPoller{Timeout: time.Second},
	})

	// turn off inside logging
	bot.Debug = false

	antiSpamBot := &Bot{
		UserName: bot.Self.UserName,
		Bot:      bot,
		BotAdm:   botAdm,
		conf:     conf,
		Storage:  storage,
	}

	return antiSpamBot, nil
}

// RegisterRoutes registers routes for tg bot in the forwarded router
func (b *Bot) RegisterRoutes(ctx context.Context) {
	b.BotAdm.Handle(admmodels.NewWord, b.InsertWord())
	b.BotAdm.Handle(admmodels.ShowWords, b.GetWords())
	b.BotAdm.Handle(admmodels.RemoveWord, b.DelWord())
}

func (b *Bot) Start() {
	b.BotAdm.Start()
}

func (b *Bot) getArg(ctx telebot.Context) (string, bool) {
	args := strings.Split(ctx.Message().Text, " ")
	if len(args) < 2 || len(args) > 2 {
		ctx.Send("после команды не введено слова или больше одного")

		return "", false
	}

	return strings.ToLower(args[1]), true
}

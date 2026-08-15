package antispambot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"

	"telegram-antispam-bot/internal/infra/config"
	admmodels "telegram-antispam-bot/internal/models/adm_models"
	models_errors_anti_spambot "telegram-antispam-bot/internal/models/errors/anti_spam_bot"
)

// antiSpamBot — интерфейс к Telegram Bot API (облегчает тестирование).
type antiSpamBot interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

// Storage — интерфейс доступа к базе данных.
type Storage interface {
	DelWordFromBadWords(word string) error
	GetListBadWords() ([]string, error)
	InsertWordToBadWords(word string) error
	GetWhitelistAuthors() ([]int64, error)
	InsertWhitelistAuthor(id int64) error
	DelWhitelistAuthor(id int64) error
	IsWhitelistAuthor(id int64) (bool, error)
	GetWhitelistTags() ([]string, error)
	InsertWhitelistTag(tag string) error
	DelWhitelistTag(tag string) error
	IsWhitelistTag(tag string) (bool, error)
}

// Bot — бот для борьбы со спамом.
type Bot struct {
	Bot      antiSpamBot
	BotAdm   *telebot.Bot
	UserName string
	Storage  Storage
	conf     *config.Config
	logger   *slog.Logger
}

// NewAntiSpamBot — конструктор AntiSpamBot. Возвращает новый экземпляр.
func NewAntiSpamBot(conf *config.Config, storage Storage, log *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(conf.BotAntiSpam.Settings.Token)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", models_errors_anti_spambot.ErrInitBot, err)
	}

	botAdm, err := telebot.NewBot(telebot.Settings{
		Token:  conf.BotAntiSpam.Settings.AdmToken,
		Poller: &telebot.LongPoller{Timeout: time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("%w:%v", models_errors_anti_spambot.ErrInitBot, err)
	}

	// отключаем внутреннее логирование библиотеки
	bot.Debug = false

	antiSpamBot := &Bot{
		UserName: bot.Self.UserName,
		Bot:      bot,
		BotAdm:   botAdm,
		conf:     conf,
		Storage:  storage,
		logger:   log,
	}

	return antiSpamBot, nil
}

// RegisterRoutes регистрирует обработчики команд админ-бота.
func (b *Bot) RegisterRoutes(ctx context.Context) {
	b.BotAdm.Handle(admmodels.NewWord, b.InsertWord())
	b.BotAdm.Handle(admmodels.ShowWords, b.GetWords())
	b.BotAdm.Handle(admmodels.RemoveWord, b.DelWord())
	b.BotAdm.Handle(admmodels.AddAuthor, b.AddAuthor())
	b.BotAdm.Handle(admmodels.RemoveAuthor, b.RemoveAuthor())
	b.BotAdm.Handle(admmodels.AddTag, b.AddTag())
	b.BotAdm.Handle(admmodels.RemoveTag, b.RemoveTag())
	b.BotAdm.Handle(admmodels.ShowWhitelist, b.ShowWhitelist())
	b.BotAdm.Handle(admmodels.Info, b.Info())
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

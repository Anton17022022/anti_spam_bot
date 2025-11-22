package antispambot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram-antispam-bot/internal/infra/config"
	models_errors_anti_spambot "telegram-antispam-bot/internal/models/errors/anti_spam_bot"
)

// antiSpamBot is API to
type antiSpamBot interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

// Bot is bot for struggle spam
type Bot struct {
	Bot      antiSpamBot
	UserName string
	conf     *config.Config
}

// NewAntiSpamBot is constructor AntiSpamBot. Return new instance.
func NewAntiSpamBot(conf *config.Config) (*Bot, error) {
	log.Println("Token %s", conf.BotAntiSpam.Settings.Token)
	
	bot, err := tgbotapi.NewBotAPI(conf.BotAntiSpam.Settings.Token)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", models_errors_anti_spambot.ErrInitBot, err)
	}

	// turn off inside logging
	bot.Debug = false

	antiSpamBot := &Bot{
		UserName: bot.Self.UserName,
		Bot:      bot,
		conf:     conf,
	}

	return antiSpamBot, nil
}

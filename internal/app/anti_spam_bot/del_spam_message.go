package antispambot

import (
	"log"
	"strings"
	"time"

	models_adds "telegram-antispam-bot/internal/models/adds"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartDelSpamMessage analyze message, and del spam.
func (b *Bot) StartDelSpamMessage() {
	u := tgbotapi.NewUpdate(b.conf.BotAntiSpam.Settings.OffsetMessageStart)
	u.Timeout = b.conf.BotAntiSpam.Settings.TimeOut

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		// TODO ограничить динамическим пулом воркеров с лимитом
		go func() {
			// check is message nil
			if update.Message != nil {
				if b.isWhiteList(update.Message) {
					return
				}

				// check is message is for del
				if b.isForDel(update.Message) {
					log.Printf("deleted spam message: chat ID: %d, user: %s, message ID: %d\n", update.Message.Chat.ID, update.Message.From.UserName, update.Message.MessageID)

					deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

					b.deleteMessageWithRetry(deleteMsg)
				}
			}
		}()
	}
}

func (b *Bot) isForDel(msg *tgbotapi.Message) bool {
	return b.containsAd(msg.Caption) || b.containsAd(msg.Text) || b.containsHyperLink(msg.Text)
}

// containsAd check if message is ad
func (b *Bot) containsAd(text string) bool {
	if text == "" {
		return false
	}

	textLower := strings.ToLower(text)

	for _, keyword := range models_adds.AdKeywords {
		if strings.Contains(textLower, keyword) {
			return true
		}
	}

	return false
}

func (b *Bot) containsHyperLink(text string) bool {
	return models_adds.HasURL(text)
}

func (b *Bot) deleteMessageWithRetry(deleteMsg tgbotapi.DeleteMessageConfig) {
	retries := b.conf.BotAntiSpam.Settings.Reties

	for i := 0; i < retries; i++ {
		if _, err := b.Bot.Request(deleteMsg); err != nil {
			log.Printf("Failed to delete message (attempt %d): %v. chat ID: %d, user: %s, message ID: %d", i+1, err.Error(), deleteMsg.ChatID, deleteMsg.ChannelUsername, deleteMsg.MessageID)

			if i == retries-1 {
				log.Println("Max retries reached, giving up.")
				return
			}

			// TODO лютый хардкод - вынести в целом реатри в раунд триппер с изменений задержки (в частности прогрессивной)
			time.Sleep(b.conf.BotAntiSpam.Settings.TimeOutBetweenRetries)

			continue
		}

		return
	}
}

func (b *Bot) isWhiteList(msg *tgbotapi.Message) bool {
	if msg.From.UserName == b.conf.BotAntiSpam.WhiteListAuthor || msg.From.UserName == "" {
		return true
	}

	if msg.Chat.UserName == b.conf.BotAntiSpam.WhiteListAuthor {
		return true
	}

	words := strings.Split(msg.Text, " ")

	for _, v := range words {
		if _, ok := b.conf.BotAntiSpam.WhiteListTags[v]; ok {
			return true
		}
	}

	return false
}

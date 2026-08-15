package antispambot

import (
	"fmt"
	"strings"
	"time"

	models_adds "telegram-antispam-bot/internal/models/adds"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"
)

// StartDelSpamMessage анализирует входящие сообщения и удаляет спам.
func (b *Bot) StartDelSpamMessage() {
	u := tgbotapi.NewUpdate(b.conf.BotAntiSpam.Settings.OffsetMessageStart)
	u.Timeout = b.conf.BotAntiSpam.Settings.TimeOut

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		// TODO: вынести обработку обновления в отдельный механизм/очередь
		go func() {
			// проверяем, что сообщение существует
			if update.Message != nil {
				senderID, ok := b.senderID(update.Message)
				if !ok {
					return
				}

				if b.isWhiteList(update.Message) {
					return
				}

				// проверяем, подлежит ли сообщение удалению
				if b.isForDel(update.Message) {
					b.logger.Info("deleted spam message",
						"chat_id", update.Message.Chat.ID,
						"user_id", senderID,
						"message_id", update.Message.MessageID,
					)

					deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

					b.deleteMessageWithRetry(deleteMsg)
				}
			}
		}()
	}
}

// isForDel определяет, является ли сообщение спамом (реклама или ссылка).
func (b *Bot) isForDel(msg *tgbotapi.Message) bool {
	return b.containsAd(msg.Caption) || b.containsAd(msg.Text) || b.containsLink(msg)
}

// containsAd проверяет, содержит ли текст рекламное/запрещённое слово.
func (b *Bot) containsAd(text string) bool {
	if text == "" {
		return false
	}

	textLower := strings.ToLower(text)

	badWords, err := b.Storage.GetListBadWords()
	if err != nil {
		b.BotAdm.Send(&telebot.User{
			ID: b.conf.BotAntiSpam.Settings.AlertUserID,
		},
			fmt.Errorf("Storage.GetListBadWords: %w", err),
		)

		return false
	}

	for _, keyword := range badWords {
		if strings.Contains(textLower, keyword) {
			return true
		}
	}

	return false
}

func (b *Bot) containsLink(msg *tgbotapi.Message) bool {
	return models_adds.HasURL(msg.Text) || models_adds.HasURL(msg.Caption) || b.isHyperLinkText(msg.Entities) || b.isHyperLinkText(msg.CaptionEntities)
}

func (b *Bot) isHyperLinkText(entities []tgbotapi.MessageEntity) bool {
	for _, v := range entities {
		if v.IsTextLink() || v.IsURL() {
			return true
		}
	}

	return false
}

func (b *Bot) deleteMessageWithRetry(deleteMsg tgbotapi.DeleteMessageConfig) {
	retries := b.conf.BotAntiSpam.Settings.Reties

	for i := 0; i < retries; i++ {
		if _, err := b.Bot.Request(deleteMsg); err != nil {
			b.logger.Warn("failed to delete message",
				"attempt", i+1,
				"error", err.Error(),
				"chat_id", deleteMsg.ChatID,
				"user", deleteMsg.ChannelUsername,
				"message_id", deleteMsg.MessageID,
			)

			if i == retries-1 {
				b.logger.Warn("max retries reached, giving up")
				return
			}

			// TODO: жёсткий хардкод — вынести ретраи в отдельный механизм (round tripper) с прогрессивной задержкой
			time.Sleep(b.conf.BotAntiSpam.Settings.TimeOutBetweenRetries)

			continue
		}

		return
	}
}

// senderID возвращает ID отправителя сообщения: пользователя или канала (от имени которого написано).
func (b *Bot) senderID(msg *tgbotapi.Message) (int64, bool) {
	if msg.From != nil {
		return msg.From.ID, true
	}

	if msg.SenderChat != nil {
		return msg.SenderChat.ID, true
	}

	return 0, false
}

func (b *Bot) isWhiteList(msg *tgbotapi.Message) bool {
	senderID, ok := b.senderID(msg)
	if !ok {
		return false
	}

	isAuthor, err := b.Storage.IsWhitelistAuthor(senderID)
	if err != nil {
		b.logger.Error("failed to check whitelist author",
			"error", err.Error(),
			"user_id", senderID,
		)
		return false
	}
	if isAuthor {
		return true
	}

	tags, err := b.Storage.GetWhitelistTags()
	if err != nil {
		b.logger.Error("failed to get whitelist tags",
			"error", err.Error(),
			"user_id", senderID,
		)
		return false
	}

	words := strings.Split(msg.Text, " ")

	for _, word := range words {
		for _, tag := range tags {
			if word == tag {
				return true
			}
		}
	}

	return false
}

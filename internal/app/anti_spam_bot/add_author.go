package antispambot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

// AddAuthor добавляет пользователя (из пересланного или reply-сообщения) в whitelist.
func (b *Bot) AddAuthor() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		id, ok := b.getTargetAuthorID(ctx)
		if !ok {
			ctx.Send("Перешлите или ответьте на сообщение пользователя.")

			return nil
		}

		isAuthor, err := b.Storage.IsWhitelistAuthor(id)
		if err != nil {
			return fmt.Errorf("Storage.IsWhitelistAuthor: %w", err)
		}
		if isAuthor {
			ctx.Send(fmt.Sprintf("Пользователь %d уже в whitelist", id))

			return nil
		}

		if err = b.Storage.InsertWhitelistAuthor(id); err != nil {
			return fmt.Errorf("Storage.InsertWhitelistAuthor: %w", err)
		}

		ctx.Send(fmt.Sprintf("Пользователь %d добавлен в whitelist", id))

		return nil
	}
}

// RemoveAuthor удаляет пользователя (из пересланного или reply-сообщения) из whitelist.
func (b *Bot) RemoveAuthor() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		id, ok := b.getTargetAuthorID(ctx)
		if !ok {
			ctx.Send("Перешлите или ответьте на сообщение пользователя.")

			return nil
		}

		if err := b.Storage.DelWhitelistAuthor(id); err != nil {
			return fmt.Errorf("Storage.DelWhitelistAuthor: %w", err)
		}

		ctx.Send(fmt.Sprintf("Пользователь %d удален из whitelist", id))

		return nil
	}
}

// getTargetAuthorID возвращает telegram id из пересланного или reply-сообщения:
// от пользователя или от канала (если сообщение пришло от имени канала).
func (b *Bot) getTargetAuthorID(ctx telebot.Context) (int64, bool) {
	msg := ctx.Message()
	if msg == nil {
		return 0, false
	}

	// пересланное сообщение: источник — пользователь или канал
	if msg.IsForwarded() {
		if msg.OriginalSender != nil {
			return msg.OriginalSender.ID, true
		}

		if msg.OriginalChat != nil {
			return msg.OriginalChat.ID, true
		}
	}

	// ответ (reply) на сообщение: автор — пользователь или канал
	if msg.ReplyTo != nil {
		if msg.ReplyTo.Sender != nil {
			return msg.ReplyTo.Sender.ID, true
		}

		if msg.ReplyTo.SenderChat != nil {
			return msg.ReplyTo.SenderChat.ID, true
		}
	}

	return 0, false
}

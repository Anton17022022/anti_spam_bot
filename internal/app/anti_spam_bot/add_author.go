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

		ids, ok := b.getTargetAuthorIDs(ctx)
		if !ok {
			ctx.Send("Перешлите или ответьте на сообщение пользователя.")

			return nil
		}

		added := 0
		for _, id := range ids {
			isAuthor, err := b.Storage.IsWhitelistAuthor(id)
			if err != nil {
				return fmt.Errorf("Storage.IsWhitelistAuthor: %w", err)
			}
			if isAuthor {
				continue
			}

			if err = b.Storage.InsertWhitelistAuthor(id); err != nil {
				return fmt.Errorf("Storage.InsertWhitelistAuthor: %w", err)
			}
			added++
		}

		if added == 0 {
			ctx.Send(fmt.Sprintf("Пользователь %v уже в whitelist", ids))

			return nil
		}

		ctx.Send(fmt.Sprintf("Пользователь %v добавлен в whitelist", ids))

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

		ids, ok := b.getTargetAuthorIDs(ctx)
		if !ok {
			ctx.Send("Перешлите или ответьте на сообщение пользователя.")

			return nil
		}

		for _, id := range ids {
			if err := b.Storage.DelWhitelistAuthor(id); err != nil {
				return fmt.Errorf("Storage.DelWhitelistAuthor: %w", err)
			}
		}

		ctx.Send(fmt.Sprintf("Пользователь %v удален из whitelist", ids))

		return nil
	}
}

// getTargetAuthorIDs возвращает telegram id из пересланного или reply-сообщения.
// Для сообщений от имени канала могут быть заполнены оба представления ID —
// пользователя (Sender/OriginalSender) и канала (SenderChat/OriginalChat), — причём в
// разных сообщениях канал фигурирует под разными формами ID. Чтобы whitelist-проверка
// гарантированно находила канал, добавляем все представления.
func (b *Bot) getTargetAuthorIDs(ctx telebot.Context) ([]int64, bool) {
	msg := ctx.Message()
	if msg == nil {
		return nil, false
	}

	ids := make([]int64, 0, 2)
	addID := func(id int64) {
		for _, existing := range ids {
			if existing == id {
				return
			}
		}
		ids = append(ids, id)
	}

	// пересланное сообщение: источник — пользователь или канал
	if msg.IsForwarded() {
		if msg.OriginalSender != nil {
			addID(msg.OriginalSender.ID)
		}

		if msg.OriginalChat != nil {
			addID(msg.OriginalChat.ID)
		}
	}

	// ответ (reply) на сообщение: автор — пользователь или канал
	if msg.ReplyTo != nil {
		if msg.ReplyTo.Sender != nil {
			addID(msg.ReplyTo.Sender.ID)
		}

		if msg.ReplyTo.SenderChat != nil {
			addID(msg.ReplyTo.SenderChat.ID)
		}
	}

	return ids, len(ids) > 0
}

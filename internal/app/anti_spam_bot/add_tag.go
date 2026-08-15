package antispambot

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

// AddTag добавляет тэг в whitelist.
func (b *Bot) AddTag() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		tagArg, ok := b.getArg(ctx)
		if !ok {
			return nil
		}

		tag := strings.ToLower(tagArg)

		isTag, err := b.Storage.IsWhitelistTag(tag)
		if err != nil {
			return fmt.Errorf("Storage.IsWhitelistTag: %w", err)
		}
		if isTag {
			ctx.Send(fmt.Sprintf("тэг %s уже в whitelist", tag))

			return nil
		}

		if err = b.Storage.InsertWhitelistTag(tag); err != nil {
			return fmt.Errorf("Storage.InsertWhitelistTag: %w", err)
		}

		ctx.Send(fmt.Sprintf("тэг %s добавлен в whitelist", tag))

		return nil
	}
}

// RemoveTag удаляет тэг из whitelist.
func (b *Bot) RemoveTag() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		tagArg, ok := b.getArg(ctx)
		if !ok {
			return nil
		}

		if err := b.Storage.DelWhitelistTag(strings.ToLower(tagArg)); err != nil {
			return fmt.Errorf("Storage.DelWhitelistTag: %w", err)
		}

		ctx.Send(fmt.Sprintf("тэг %s удален из whitelist", tagArg))

		return nil
	}
}

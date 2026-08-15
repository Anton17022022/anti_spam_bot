package antispambot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

// ShowWhitelist показывает текущие списки whitelist: авторов и тэгов.
func (b *Bot) ShowWhitelist() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		authors, err := b.Storage.GetWhitelistAuthors()
		if err != nil {
			return fmt.Errorf("Storage.GetWhitelistAuthors: %w", err)
		}

		tags, err := b.Storage.GetWhitelistTags()
		if err != nil {
			return fmt.Errorf("Storage.GetWhitelistTags: %w", err)
		}

		ctx.Send(fmt.Sprintf("Авторы: %v\nТэги: %v", authors, tags))

		return nil
	}
}

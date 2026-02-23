package antispambot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

// GetWords ..
func (b *Bot) GetWords() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		words, err := b.Storage.GetListBadWords()
		if err != nil {
			return fmt.Errorf("Storage.GetListBadWords: %w", err)
		}

		ctx.Send(fmt.Sprintf("Список: %v", words))

		return nil
	}
}

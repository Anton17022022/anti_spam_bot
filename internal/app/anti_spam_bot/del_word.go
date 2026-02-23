package antispambot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

// DelWord ..
func (b *Bot) DelWord() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		wordArg, ok := b.getArg(ctx)
		if !ok {
			return nil
		}

		err := b.Storage.DelWordFromBadWords(wordArg)
		if err != nil {
			return fmt.Errorf("Storage.DelWordFromBadWords: %w", err)
		}

		ctx.Send(fmt.Sprintf("слово %s удалено", wordArg))

		return nil
	}
}

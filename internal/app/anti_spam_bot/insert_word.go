package antispambot

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

// InsertWord ..
func (b *Bot) InsertWord() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		if !b.Auth(ctx.Sender().ID) {
			ctx.Send("Чеши от сэда")

			return nil
		}

		wordArg, ok := b.getArg(ctx)
		if !ok {
			return nil
		}

		words, err := b.Storage.GetListBadWords()
		if err != nil {
			return fmt.Errorf("Storage.GetListBadWords: %w", err)
		}

		for _, word := range words {
			if word == wordArg {
				ctx.Send(fmt.Sprintf("слово %s есть в списке", wordArg))

				return nil
			}
		}

		err = b.Storage.InsertWordToBadWords(strings.ToLower(wordArg))
		if err != nil {
			return fmt.Errorf("Storage.InsertWordToBadWords: %w", err)
		}

		ctx.Send(fmt.Sprintf("Слово %s добавлено в список запрещенный слов.", wordArg))

		return nil
	}
}

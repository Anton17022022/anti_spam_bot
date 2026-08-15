package antispambot

import (
	"gopkg.in/telebot.v3"
)

// infoText — описание функционала всех команд админ-бота.
const infoText = `Команды админ-бота:

Слова (запрещённые):
/new_word <слово> — добавить слово в список запрещённых
/remove_word <слово> — удалить слово из списка
/show_words — показать текущий список запрещённых слов

Whitelist (игнорируется антиспамом):
/add_author — ответьте или перешлите сообщение пользователя, чтобы добавить его
/remove_author — ответьте или перешлите сообщение пользователя, чтобы убрать его
/add_tag <тэг> — добавить тэг в whitelist
/remove_tag <тэг> — удалить тэг из whitelist
/show_whitelist — показать авторов и тэги из whitelist

Прочее:
/info — показать это описание

Команды изменения списков доступны только пользователям из whitelist авторов.`

// Info возвращает обработчик команды показа справки по командам бота.
func (b *Bot) Info() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		ctx.Send(infoText)

		return nil
	}
}

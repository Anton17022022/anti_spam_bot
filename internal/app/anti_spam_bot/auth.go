package antispambot

// Auth ..
func (b *Bot) Auth(authorID int64) bool {
	return authorID == b.conf.BotAntiSpam.WhiteListAuthor
}

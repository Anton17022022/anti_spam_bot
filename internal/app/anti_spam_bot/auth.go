package antispambot

// Auth ..
func (b *Bot) Auth(authorID int64) bool {
	for _, whAutor := range b.conf.BotAntiSpam.WhiteListAuthor {
		if authorID == whAutor {
			return true
		}
	}

	return false
}

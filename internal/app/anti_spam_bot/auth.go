package antispambot

// Auth проверяет, входит ли автор в whitelist.
func (b *Bot) Auth(authorID int64) bool {
	isAuthor, err := b.Storage.IsWhitelistAuthor(authorID)
	if err != nil {
		b.logger.Error("failed to check whitelist author",
			"error", err.Error(),
			"author_id", authorID,
		)
		return false
	}

	return isAuthor
}

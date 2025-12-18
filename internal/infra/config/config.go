package config

import (
	"os"
	"time"
)

// Config is config app
type Config struct {
	BotAntiSpam botAntiSpam
}

type botAntiSpam struct {
	Settings        settings
	WhiteListTags   map[string]struct{}
	WhiteListAuthor int64
}

type settings struct {
	Token                 string
	OffsetMessageStart    int
	TimeOut               int // for long request to interrapt
	Reties                int // when retries to del message
	TimeOutBetweenRetries time.Duration
}

// NewConfig return config app instance
func NewConfig() (*Config, error) {
	wlTags := map[string]struct{}{
		"@prolann": {},
		"@Prolann": {},
	}

	return &Config{
		BotAntiSpam: botAntiSpam{
			Settings: settings{
				Token:                 os.Getenv("antispam_bot_token"),
				OffsetMessageStart:    0,
				TimeOut:               60,
				Reties:                3,
				TimeOutBetweenRetries: 10 * time.Second,
			},
			WhiteListTags:   wlTags,
			WhiteListAuthor: 136817688,
		}}, nil
}

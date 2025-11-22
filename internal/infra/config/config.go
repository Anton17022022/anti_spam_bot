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
	Settings      settings
	WhiteListTags map[string]struct{}
}

type settings struct {
	Token                 string
	OffsetMessageStart    int
	TimeOut               int // for long request to interrapt
	Reties                int // when retries to del message
	TimeOutBetweenRetries time.Duration
}

// NewConfig return config app instance
func NewConfig() *Config {
	wlTags := map[string]struct{}{
		"@prolann": {},
		"@Prolann": {},
	}

	return &Config{
		BotAntiSpam: botAntiSpam{
			Settings: settings{
				Token:                 os.Getenv("bot_token"),
				OffsetMessageStart:    0,
				TimeOut:               60,
				Reties:                3,
				TimeOutBetweenRetries: 10 * time.Second,
			},
			WhiteListTags: wlTags,
		}}
}

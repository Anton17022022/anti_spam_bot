package main

import (
	"context"

	"telegram-antispam-bot/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		panic(err.Error())
	}

	a.ListenAndServe()
}

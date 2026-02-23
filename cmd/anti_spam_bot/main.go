package main

import (
	"context"
	"log"
	"telegram-antispam-bot/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		panic(err.Error())
	}

	log.Println("app inited")
	a.ListenAndServe()
}

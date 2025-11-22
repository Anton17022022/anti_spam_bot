package main

import (
	"log"
	"telegram-antispam-bot/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		panic(err.Error())
	}

	log.Println("app inited")
	a.ListenAndServe()
}

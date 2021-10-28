package main

import (
	"fmt"
	"next-german-words/app/tg"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	tgApiKey := os.Getenv("TELEGRAM_API_KEY")
	telegram := tg.NewTelegramClient(tgApiKey)
	go telegram.Start()
	fmt.Println("Press CTRL + C to stop programm")
	select {
	case <- sigCh:
		fmt.Println("Shutting down")
		os.Exit(0)
	}
}

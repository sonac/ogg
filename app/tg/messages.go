package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"next-german-words/app/store/models"
	"strings"
)

// Responds to different messages

func (tg *Telegram) start(chatId int64) {
	usr := &models.User{TelegramChatId: chatId}
	alreadyExists, err := tg.DB.InsertUser(usr)
	if err != nil {
		tg.sendMessage("Internal error occured", chatId)
		log.Println(err)
		return
	}
	if alreadyExists {
		tg.sendMessage("You already registered at snitch", chatId)
		return
	}
	tg.sendMessage("Welcome!", chatId)
}

func (tg *Telegram) addWord(msg *tgbotapi.Message) {
	word := strings.Split(msg.Text, " ")[1]
	translatedWord, err := tg.Tr.TranslateWord(word)
	if err != nil {
		tg.sendMessage("internal error occurred", msg.Chat.ID)
		log.Println(err)
		return
	}
	err = tg.DB.InsertWord(translatedWord)
	if err != nil {
		tg.sendMessage("internal error occurred", msg.Chat.ID)
		log.Println(err)
		return
	}
	err = tg.Tr.RefreshWords()
	if err != nil {
		tg.sendMessage("internal error occurred", msg.Chat.ID)
		log.Println(err)
		return
	}
}
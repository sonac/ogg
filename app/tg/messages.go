package tg

import (
	"fmt"
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
	tg.sendListManageButton(chatId)
}

func (tg *Telegram) addWord(msg *tgbotapi.Message) {
	word := strings.Split(msg.Text, " ")[1]
	translatedWord, err := tg.Tr.TranslateWord(word)
	if err != nil {
		tg.sendMessage("internal error occurred", msg.Chat.ID)
		log.Println(err)
		return
	}
	if translatedWord == nil {
		tg.sendMessage("there is no such noun in english", msg.Chat.ID)
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
	tg.sendMessage(fmt.Sprintf("%s now in the database", translatedWord.German), msg.Chat.ID)
}

func (tg *Telegram) sendRandomWord(chatId int64) {
	usr, err := tg.DB.GetUserByChatId(chatId)
	word := tg.Tr.GetRandomWord(&usr.CurStreakWords)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
	usr.CurWord = *word
	err = tg.DB.UpdateUser(usr)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("der", "der"),
			tgbotapi.NewInlineKeyboardButtonData("die", "die"),
			tgbotapi.NewInlineKeyboardButtonData("das", "das"),
			),
		)
	msg := tgbotapi.NewMessage(chatId, word.German)
	msg.ReplyMarkup = keyboard
	_, err = tg.Bot.Send(msg)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
}

func (tg *Telegram) sendBestStreak(chatId int64) {
	usr, err := tg.DB.GetUserByChatId(chatId)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
	tg.sendMessage(fmt.Sprintf("Your best streak is %d", usr.BestStreak), chatId)
}

func (tg *Telegram) checkAnswer(cb *tgbotapi.CallbackQuery) {
	chatId := cb.Message.Chat.ID
	ans := cb.Data
	usr, err := tg.DB.GetUserByChatId(chatId)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
	correctAns := genToArticle(usr.CurWord.Gen)
	if correctAns == ans {
		usr.CurStreak += 1
		usr.CurStreakWords = append(usr.CurStreakWords, usr.CurWord.German)
		err = tg.DB.UpdateUser(usr)
		if err != nil {
			log.Println(err)
			tg.sendMessage("internal error occured", chatId)
			return
		}
		tg.sendMessage("Correct!", chatId)
	} else {
		curStreak := usr.CurStreak
		usr.CurStreak = 0
		usr.CurStreakWords = nil
		err = tg.DB.UpdateUser(usr)
		if err != nil {
			log.Println(err)
			tg.sendMessage("internal error occured", chatId)
			return
		}
		msg := fmt.Sprintf("Wrong! Correct answer is %s \nYou've got %d words in a row", correctAns, curStreak)
		tg.sendMessage(msg, chatId)
	}
	tg.sendRandomWord(chatId)
}

func articleToGen(article string) string {
	switch article {
	case "der":
		return "m"
	case "die":
		return "f"
	case "das":
		return "n"
	}
	return ""
}

func genToArticle(gen string) string {
	switch gen {
	case "m":
		return "der"
	case "f":
		return "die"
	case "n":
		return "das"
	}
	return ""
}
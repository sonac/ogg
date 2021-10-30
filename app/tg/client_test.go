package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"next-german-words/app/store/models"
	mocks "next-german-words/mocks/app/tg"
	"testing"
)

var (
	telegram Telegram
	mockedMongo *mocks.Database
	mockedTranslator *mocks.Translator
)

func init() {
	mockedMongo = &mocks.Database{}
	mockedTranslator = &mocks.Translator{}
	telegram = Telegram{DB: mockedMongo, Tr: mockedTranslator}
}

func TestSendMessage(t *testing.T) {
	t.Run("Testing simple message", func(t *testing.T) {
		mockedBot := &mocks.TelegramBot{}
		telegram.Bot = mockedBot
		msgText := "some message"
		msg := tgbotapi.NewMessage(123, msgText)
		mockedBot.On("Send", msg).Return(tgbotapi.Message{}, nil)
		telegram.sendMessage(msgText, 123)
		mockedBot.AssertExpectations(t)
		mockedBot.AssertCalled(t, "Send", msg)
	})
}

func TestCheckAnswer(t *testing.T) {
	t.Run("Testing whether the answer is being checked correctly", func(t *testing.T) {
		chatId := int64(123)
		mockedBot := &mocks.TelegramBot{}
		telegram.Bot = mockedBot
		chat := tgbotapi.Chat{
			ID:                  chatId,
			Type:                "",
			Title:               "",
			UserName:            "",
			FirstName:           "",
			LastName:            "",
			AllMembersAreAdmins: false,
			Photo:               nil,
			Description:         "",
			InviteLink:          "",
		}
		msg := tgbotapi.Message{}
		msg.Chat = &chat
		upd := tgbotapi.CallbackQuery{}
		upd.Message = &msg
		upd.Data = "die"
		word := models.Word{
			German:  "Wand",
			English: "Wall",
			Gen:     "f",
			Meaning: "none",
		}
		usr := models.User{
			TelegramChatId: 123,
			BestStreak:     1,
			CurWord:        word,
			CurStreak:      2,
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("der", "der"),
				tgbotapi.NewInlineKeyboardButtonData("die", "die"),
				tgbotapi.NewInlineKeyboardButtonData("das", "das"),
			),
		)
		updatedUser := usr
		updatedUser.CurStreak += 1
		updatedUser.CurStreakWords = []string{"Wand"}
		mockedMongo.On("GetUserByChatId", chatId).Return(&usr, nil)
		mockedMongo.On("UpdateUser", &updatedUser).Return(nil)
		correctMsg := tgbotapi.NewMessage(chatId, "Correct!")
		wordMsg := tgbotapi.NewMessage(chatId, "Wand")
		wordMsg.ReplyMarkup = keyboard
		mockedBot.On("Send", wordMsg).Return(tgbotapi.Message{}, nil)
		mockedBot.On("Send", correctMsg).Return(tgbotapi.Message{}, nil)
		mockedTranslator.On("GetRandomWord", &updatedUser.CurStreakWords).Return(&word)
		telegram.checkAnswer(&upd)
		mockedBot.AssertCalled(t, "Send", correctMsg)
		mockedBot.AssertCalled(t, "Send", wordMsg)
		if usr.CurStreak != 3 {
			t.Error("Streak should've been increased")
		}
	})
}
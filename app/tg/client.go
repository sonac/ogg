package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"next-german-words/app/store"
	"next-german-words/app/store/models"
	"next-german-words/app/translator"
	"next-german-words/app/utils"
	"os"
)

type TelegramBot interface {
	GetUpdatesChan(tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error)
	AnswerCallbackQuery(tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Database interface {
	Connect()
	InsertWord(word *models.Word) error
	FindWord(german string) (*models.Word, error)
	GetWords() ([]*models.Word, error)
	InsertUser(user *models.User) (bool, error)
	GetUserByChatId (chatId int64) (*models.User, error)
	UpdateUser(user *models.User) error
}

type Translator interface {
	RefreshWords() error
	TranslateWord(word string) (*models.Word, error)
	GetRandomWord(wordsToFilter *[]string) *models.Word
}

type Telegram struct {
	Bot TelegramBot
	DB Database
	Tr Translator
	updateConfig tgbotapi.UpdateConfig
}

func NewTelegramClient(apiKey string) *Telegram {
	tg := Telegram{}
	tg.Init(apiKey)
	return &tg
}

func (tg *Telegram) Init(apiKey string) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Println("")
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		bot.Debug = true
	}
	tg.updateConfig = tgbotapi.NewUpdate(0)
	tg.updateConfig.Timeout = 60
	tg.Bot = bot
	tg.DB = &store.Mongo{}
	tg.DB.Connect()
	tg.Tr = translator.NewTranslator(utils.Getenv("YANDEX_API_KEY", ""))
}

func (tg *Telegram) Start() {
	tg.readUpdates(tg.updateConfig)
}

func (tg *Telegram) readUpdates(updateConfig tgbotapi.UpdateConfig) {
	updates, err := tg.Bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Println("[ERROR] There were an error during tg initialization")
		log.Fatalln(err)
	}
	for upd := range updates {
		if upd.CallbackQuery != nil {
			tg.replyToCallback(upd.CallbackQuery)
		}
		if upd.Message != nil {
			tg.replyToMessage(upd.Message)
		}
	}
}

func (tg *Telegram) replyToMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() {
		tg.replyToCommand(msg)
		return
	}
	switch msg.Text {
	case "start":
		tg.start(msg.Chat.ID)
	case "Next Word":
		tg.sendRandomWord(msg.Chat.ID)
	case "Best Streak":
		tg.sendBestStreak(msg.Chat.ID)
	case "Info":
		tg.sendMessage("He-he, nothing here yet", msg.Chat.ID)
	default:
		tg.sendMessage("I don't understand", msg.Chat.ID)
	}
}

func (tg *Telegram) replyToCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		tg.start(msg.Chat.ID)
	case "add":
		tg.addWord(msg)
	case "restart":
		tg.sendListManageButton(msg.Chat.ID)
	}
}

func (tg *Telegram) replyToCallback(cb *tgbotapi.CallbackQuery) {
	tg.checkAnswer(cb)
}

func (tg *Telegram) sendMessage(msg string, chatId int64) {
	message := tgbotapi.NewMessage(chatId, msg)
	_, err := tg.Bot.Send(message)
	if err != nil {
		log.Printf("[ERROR] Occured during sending message to tg chat %d", chatId)
		log.Fatalln(err)
	}
}

func (tg *Telegram) sendListManageButton(chatId int64) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Next Word"),
			tgbotapi.NewKeyboardButton("Best Streak"),
			tgbotapi.NewKeyboardButton("Info"),
		),
	)
	message := tgbotapi.NewMessage(chatId, "Added list view keyboard")
	message.ReplyMarkup = keyboard
	_, err := tg.Bot.Send(message)
	if err != nil {
		tg.sendMessage("internal error occurred", chatId)
		log.Println(err)
		return
	}
}

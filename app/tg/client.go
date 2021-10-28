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
}

type Translator interface {
	RefreshWords() error
	TranslateWord(word string) (*models.Word, error)
	GetRandomWord() *models.Word
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
		log.Println("Received a message: ", upd.Message.Text)
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
	default:
		tg.sendMessage("I don't understand", msg.Chat.ID)
	}
}

func (tg *Telegram) replyToCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "add":
		tg.addWord(msg)
	}
}

func (tg *Telegram) sendMessage(msg string, chatId int64) {
	message := tgbotapi.NewMessage(chatId, msg)
	_, err := tg.Bot.Send(message)
	if err != nil {
		log.Printf("[ERROR] Occured during sending message to tg chat %d", chatId)
		log.Fatalln(err)
	}
}

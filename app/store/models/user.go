package models

type User struct {
	TelegramChatId int64 `bson:"telegram_chat_id"`
	BestStreak int64 `bson:"best_streak"`
}
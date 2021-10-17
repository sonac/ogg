package models

type User struct {
	TelegramChat int64 `bson:"telegram_chat"`
	BestStreak int64 `bson:"best_streak"`
}
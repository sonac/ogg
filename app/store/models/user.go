package models

type User struct {
	TelegramChatId int64 `bson:"telegram_chat_id"`
	BestStreak int64 `bson:"best_streak"`
	CurWord Word `bson:"cur_word"`
	CurStreakWords []string `bson:"cur_streak_words"`
	CurStreak int64 `bson:"cur_streak"`
}
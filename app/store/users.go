package store

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"next-german-words/app/store/models"
)

func (m *Mongo) InsertUser(user *models.User) (bool, error) {
	exists, err := m.isUserExists(user.TelegramChatId)
	if err != nil {
		return exists, err
	}
	if exists {
		return true, nil
	}
	_, err = m.UserCollection.InsertOne(ctx, user)
	return false, err
}

func (m *Mongo) GetUserByChatId (chatId int64) (*models.User, error) {
	fltr := bson.M{"telegram_chat_id": chatId}
	res := m.UserCollection.FindOne(ctx, fltr)
	var usr models.User
	err := res.Decode(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (m *Mongo) UpdateUser(user *models.User) error {
	exists, err := m.isUserExists(user.TelegramChatId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("trying to update user that doesn't exist")
	}
	_, err = m.UserCollection.UpdateOne(ctx,
		bson.M{"telegram_chat_id": user.TelegramChatId},
		bson.M{"$set": bson.M{
			"best_streak":	user.BestStreak,
		}},
	)
	return err
}

func (m *Mongo) isUserExists(chatId int64) (bool, error) {
	cnt, err := m.UserCollection.CountDocuments(ctx, bson.D{primitive.E{Key: "telegram_chat_id", Value: chatId}})
	if err != nil {
		return false, err
	}
	return cnt >= 1, nil
}
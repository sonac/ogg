package store

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"next-german-words/app/store/models"
)

func (m *Mongo) FindWord(german string) (*models.Word, error) {
	fltr := bson.D{primitive.E{Key: "german", Value: german}}
	res := m.WordCollection.FindOne(ctx, fltr)
	var word *models.Word
	err := res.Decode(word)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (m *Mongo) GetWords() ([]*models.Word, error) {
	res, err := m.WordCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var words []*models.Word
	for res.Next(ctx) {
		var w models.Word
		err = res.Decode(&w)
		if err != nil {
			return nil, err
		}
		words = append(words, &w)
	}
	return words, nil
}

func (m *Mongo) InsertWord(word *models.Word) error {
	if b, err := m.isWordExists(word.German); err != nil || b {
		return err
	}
	_, err := m.WordCollection.InsertOne(ctx, word)
	return err
}

func (m *Mongo) isWordExists(german string) (bool, error) {
	cnt, err := m.WordCollection.CountDocuments(ctx, bson.D{primitive.E{Key: "german", Value: german}})
	if err != nil {
		return false, err
	}
	return cnt >= 1, nil
}

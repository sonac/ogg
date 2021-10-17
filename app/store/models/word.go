package models

type Word struct {
	German string `bson:"german"`
	English string `bson:"english"`
	Gen string `bson:"gen"`
	Meaning string `bson:"meaning"`
}

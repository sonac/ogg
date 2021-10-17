package translator

import "next-german-words/app/store"

type Client struct {
	YandexApiKey string
	Database *store.Mongo
}

func NewTranslator(yandexApiKey string) *Client {
	db := store.NewDatabase()
	return &Client{YandexApiKey: yandexApiKey, Database: db}
}

func (c *Client) GetWord() {

}


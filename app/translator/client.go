package translator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"next-german-words/app/store"
	"next-german-words/app/store/models"
	"time"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	YandexApiKey string
	Database     *store.Mongo
	Words        []*models.Word
}

type Translation struct {
	Text string `json:"text,omitempty"`
	Pos  string `json:"pos,omitempty"`
	Gen  string `json:"gen,omitempty"`
}

type WordDef struct {
	Text string      `json:"text,omitempty"`
	Pos  string      `json:"pos,omitempty"`
	Ts   string      `json:"ts,omitempty"`
	Tr   Translation `json:"tr"`
}

type TranslateResp struct {
	Def WordDef `json:"def"`
}

var (
	httpClient HTTPClient
)

func init() {
	httpClient = &http.Client{}
}

func NewTranslator(yandexApiKey string) *Client {
	db := store.NewDatabase()
	words, err := db.GetWords()
	if err != nil {
		log.Fatalf("failed to initialize translator, %s", err)
	}
	return &Client{YandexApiKey: yandexApiKey, Database: db, Words: words}
}

func (c *Client) GetRandomWord() *models.Word {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return c.Words[rand.Intn(len(c.Words))]
}

func (c *Client) TranslateWord(word string) (*models.Word, error) {
	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-de&text=%s", c.YandexApiKey, word)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp TranslateResp
	rawResp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	w := models.Word{
		German:  resp.Def.Tr.Text,
		English: word,
		Gen:     resp.Def.Tr.Gen,
		Meaning: "",
	}
	return &w, nil
}

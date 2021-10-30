package translator

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"next-german-words/app/store/models"
	"testing"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockHTTP struct {
	MockDo MockDoType
}

func (m *MockHTTP) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

type MockGetWordsType func() ([]*models.Word, error)

type MockDatabase struct {
	MockGetWords MockGetWordsType
}

func (m *MockDatabase) GetWords() ([]*models.Word, error) {
	return nil, nil
}

func TestClient_TranslateWord(t *testing.T) {
	db := &MockDatabase{}
	translator := Client{database: db}
	t.Run("Testing retreiving price", func(t *testing.T) {
		jsBody :=
			`{
		   "head":{
			  
		   },
		   "def":[
			  {
				 "text":"castle",
				 "pos":"noun",
				 "ts":"kɑːsl",
				 "tr":[
					{
					   "text":"Schloss",
					   "pos":"noun",
					   "gen":"n",
					   "fr":10,
					   "syn":[
						  {
							 "text":"Burg",
							 "pos":"noun",
							 "gen":"f",
							 "fr":10
						  },
						  {
							 "text":"Castle",
							 "pos":"noun",
							 "gen":"n",
							 "fr":10
						  }
					   ],
					   "mean":[
						  {
							 "text":"palace"
						  },
						  {
							 "text":"fortress"
						  }
					   ]
					}
				 ]
			  }
		   ]
		}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsBody)))
		resp := &http.Response{StatusCode: 200, Body: body}
		httpClient = &MockHTTP{
			MockDo: func(*http.Request) (*http.Response, error) {
				return resp, nil
			},
		}
		res, err := translator.TranslateWord("castle")
		if err != nil || res.German != "Schloss"{
			t.Errorf("TestClient_TranslateWord failed, %s", err)
		}
	})
	t.Run("Empty resp", func(t *testing.T) {
		jsBody := `{
			"head": {},
			"def": []
		}`
		body := ioutil.NopCloser(bytes.NewReader([]byte(jsBody)))
		resp := &http.Response{StatusCode: 200, Body: body}
		httpClient = &MockHTTP{
			MockDo: func(*http.Request) (*http.Response, error) {
				return resp, nil
			},
		}
		res, err := translator.TranslateWord("castle")
		if err != nil || res != nil {
			t.Errorf("we need to successfully return nil, %s", err)
		}
	})
}

func TestClient_filteredWords(t *testing.T) {
	words := []*models.Word{
		{"Wand", "Wall", "f", ""},
		{"Tisch", "Table", "m", ""},
	}
	translator := Client{words: words}
	fWords := translator.filteredWords(&[]string{"Wand"})
	if fWords[0].German != "Tisch" {
		t.Errorf("wrong result of filtered words")
	}
}
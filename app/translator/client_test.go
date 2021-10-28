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
	translator := NewTranslator("fooApiKey")
	t.Run("Testing retreiving price", func(t *testing.T) {
		res, err := translator.TranslateWord("castle")
		if err != nil || res.German != "Schloss"{
			t.Errorf("TestClient_TranslateWord failed, %s", err)
		}
	})
}

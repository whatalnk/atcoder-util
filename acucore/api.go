package acucore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

// Submission JSON types
type Submission struct {
	ID               int
	EpochSecond      float64 `json:"epoch_second"`
	Point, Length    float32
	ExecutionTime    float32 `json:"execution_time"`
	ProblemID        string  `json:"problem_id"`
	ContestID        string  `json:"contest_id"`
	UserID           string  `json:"user_id"`
	Language, Result string
}

// FetchSubmissions fetch submissions
func FetchSubmissions(user string) map[int]Submission {
	url := "https://kenkoooo.com/atcoder/atcoder-api/results?user=" + user
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Decode JSON
	decoder := json.NewDecoder(resp.Body)
	var data []Submission
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[int]Submission)

	for _, v := range data {
		m[v.ID] = v
	}

	return m
}

func fetchCode(contestID string, id int) string {
	url := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%d", contestID, id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	submissionCode := GetElementByID(doc, "submission-code")
	ol := submissionCode.FirstChild
	code := ol.Data
	return code
}

// Update update files
func Update(targetDir string, submissions map[int]Submission) {
	for _, v := range submissions {
		if v.Result != "AC" {
			continue
		}
		lang := NormLang(v.Language)
		ext := LangToExt(lang)
		fileName := fmt.Sprintf("%d%s", v.ID, ext)
		filePath := filepath.Join(targetDir, lang, v.ContestID, v.ProblemID, fileName)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(filePath), 0777)
			if err != nil {
				log.Fatal(err)
			}
			code := fetchCode(v.ContestID, v.ID)
			// fmt.Println(code)
			save(filePath, code)
		}
		// for debug
		break
	}
}

func save(filePath string, code string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

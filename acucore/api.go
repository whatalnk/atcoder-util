package acucore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/html"
)

var problems map[ProblemKey]string
var location *time.Location

// Problem correspond to data from information API
type Problem struct {
	ID        string
	ContestID string `json:"contest_id"`
	Title     string
}

// ProblemKey map key to store problems
type ProblemKey struct {
	ID, ContestID string
}

// FetchProblems fetch problems using information API
func FetchProblems() map[ProblemKey]string {
	url := "https://kenkoooo.com/atcoder/resources/problems.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Decode JSON
	decoder := json.NewDecoder(resp.Body)
	var data []Problem
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[ProblemKey]string)

	for _, v := range data {
		m[ProblemKey{v.ContestID, v.ID}] = v.Title
	}

	return m
}

// Submission correspond to data from Submission API
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

func codeMetaData(contestID string, problemID string, language string, epochSecond float64, id int) string {
	lang := NormLang(language)
	kwComment := KwComment(lang)
	ret := fmt.Sprintf("%s Contest ID: %s\n", kwComment, contestID)
	ret += fmt.Sprintf("%s Problem ID: %s ( https://atcoder.jp/contests/%s/tasks/%s )\n", kwComment, problemID, contestID, problemID)
	ret += fmt.Sprintf("%s Title: %s\n", kwComment, problems[ProblemKey{contestID, problemID}])
	ret += fmt.Sprintf("%s Language: %s\n", kwComment, language)
	ret += fmt.Sprintf("%s Submitted: %s ( https://atcoder.jp/contests/%s/submissions/%d ) \n\n", kwComment, time.Unix(int64(epochSecond), 0).In(location), contestID, id)
	return ret
}

// Update update files
func Update(targetDir string, submissions map[int]Submission) {
	problems = FetchProblems()
	var err error
	location, err = time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}
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
			code := codeMetaData(v.ContestID, v.ProblemID, v.Language, v.EpochSecond, v.ID)
			code += fetchCode(v.ContestID, v.ID)
			save(filePath, code)
		}
	}
}

func save(filePath string, code string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(code)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved: %s", filePath)
}

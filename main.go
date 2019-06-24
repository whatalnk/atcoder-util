package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func fetchSubmissions(user string) map[int]Submission {
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

func normLang(l string) string {
	var ret string
	switch {
	case strings.Contains(l, "Ruby"):
		ret = "ruby"
	case strings.Contains(l, "Python"):
		ret = "python"
	case strings.Contains(l, "C++"):
		ret = "cpp"
	case strings.Contains(l, "Go"):
		ret = "go"
	case strings.Contains(l, "Rust"):
		ret = "rust"
	default:
		ret = "other"
	}
	return ret
}

func langToExt(l string) string {
	var ret string
	switch l {
	case "ruby":
		ret = ".rb"
	case "python":
		ret = ".py"
	case "cpp":
		ret = ".cpp"
	case "go":
		ret = ".go"
	case "rust":
		ret = ".rs"
	default:
		ret = ".txt"
	}
	return ret
}

func checkAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false

}

func checkID(n *html.Node, key string) bool {
	if n.Type == html.ElementNode {
		v, ok := checkAttribute(n, "id")
		if ok && v == key {
			return true
		}
	}
	return false
}

func walk(n *html.Node, key string) *html.Node {
	if checkID(n, key) {
		return n
	}

	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		res := walk(nn, key)
		if res != nil {
			return res
		}
	}
	return nil
}

func getElementByID(n *html.Node, key string) *html.Node {
	return walk(n, key)
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
	submissionCode := getElementByID(doc, "submission-code")
	ol := submissionCode.FirstChild
	code := ol.Data
	return code
}

func update(targetDir string, submissions map[int]Submission) {
	for _, v := range submissions {
		lang := normLang(v.Language)
		ext := langToExt(lang)
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

	f.WriteString(code)
}

func main() {
	// Get API
	user := "wtnk0812"
	targetDir := "atcoder"
	submissions := fetchSubmissions(user)
	update(targetDir, submissions)
	// log.Println(submissions[5429957])
}

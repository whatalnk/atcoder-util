package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// FileInfo types for files to be saved
type FileInfo struct {
	EpochSecond                    float64
	ProblemID, ContestID, Language string
}

func main() {
	// Get API
	user := "wtnk0812"
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

	m := make(map[int]FileInfo)

	for _, v := range data {
		m[int(v.ID)] = FileInfo{
			EpochSecond: v.EpochSecond,
			ProblemID:   v.ProblemID,
			ContestID:   v.ContestID,
			Language:    v.Language,
		}
	}

	// save gob to local file
	f, err := os.Create("submissions.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Save: Done!")

	loaded := make(map[int]FileInfo)
	f, err = os.Open("submissions.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	err = dec.Decode(&loaded)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Load: Done!")
	log.Println(loaded[5429957])
}

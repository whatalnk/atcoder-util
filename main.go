package main

import "github.com/whatalnk/atcoder-util/acucore"

func main() {
	// Get API
	user := "wtnk0812"
	targetDir := "atcoder"
	submissions := acucore.FetchSubmissions(user)
	acucore.Update(targetDir, submissions)
	// log.Println(submissions[5429957])
}

package main

import "github.com/whatalnk/atcoder-util/acucore"

func main() {
	user := "wtnk0812"
	targetDir := "atcoder"
	submissions := acucore.FetchSubmissions(user)
	acucore.Update(targetDir, submissions)
}

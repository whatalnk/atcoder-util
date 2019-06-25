package acucore

import (
	"strings"
)

// NormLang normalize languages
func NormLang(l string) string {
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

// LangToExt language to extention
func LangToExt(l string) string {
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

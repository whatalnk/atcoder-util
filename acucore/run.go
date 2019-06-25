package acucore

// Config config
type Config struct {
	User, TargetDir string
}

func run(c Config) {
	submissions := FetchSubmissions(c.User)
	Update(c.TargetDir, submissions)
}

// Run run
func Run(c Config) {
	run(c)
}

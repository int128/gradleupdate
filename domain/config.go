package domain

type Config struct {
	GitHubToken string
	CSRFKey     string // 32 bytes string
}

func (c Config) IsValid() bool {
	return c.GitHubToken != "" && c.CSRFKey != ""
}

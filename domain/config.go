package domain

type Config struct {
	GitHubToken string
}

func (c Config) IsZero() bool {
	return c.GitHubToken == ""
}

package config

type Credentials struct {
	GitHubToken string // (required)
	CSRFKey     []byte // (required) 32 bytes
}

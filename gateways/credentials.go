package gateways

import (
	"context"
	"encoding/base64"
	"os"
	"sync"

	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

/*
NewCredentials returns an implementation of gateways.Credentials.

It resolves credentials from the following resources in order:

	1. In-memory cache
	2. Environment variables
	3. Cloud Datastore

You need to create the following entity on Datastore:

	* Kind = Credentials
	* Key (string) = DEFAULT
	* Property GitHubToken (string) = your GitHub token
	* Property CSRFKey (string) = base64 encoded string of 32 bytes key

You can give credentials by the environment variables instead of Datastore.

	GITHUB_TOKEN=token
	CSRF_KEY=base64key

You can generate `CSRFKey` by the following command:

	dd if=/dev/random bs=32 count=1 | base64

*/
func NewCredentials(logger gateways.Logger) gateways.Credentials {
	return &credentialsCache{
		Base: &credentialsIfEnv{
			Logger: logger,
			Base:   &credentialsData{},
		},
	}
}

type credentialsCache struct {
	Base gateways.Credentials
	l    sync.Mutex
	v    *config.Credentials
}

func (r *credentialsCache) Get(ctx context.Context) (*config.Credentials, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if r.v != nil {
		return r.v, nil
	}
	v, err := r.Base.Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting credentials")
	}
	r.v = v
	return r.v, nil
}

type credentialsIfEnv struct {
	Base   gateways.Credentials
	Logger gateways.Logger
}

func (r *credentialsIfEnv) Get(ctx context.Context) (*config.Credentials, error) {
	if os.Getenv("GITHUB_TOKEN") == "" || os.Getenv("CSRF_KEY") == "" {
		return r.Base.Get(ctx)
	}

	r.Logger.Infof(ctx, "fallback to credentials from environment variables")
	csrfKey, err := base64.StdEncoding.DecodeString(os.Getenv("CSRF_KEY"))
	if err != nil {
		return nil, errors.Wrapf(err, "error while decoding base64")
	}
	return &config.Credentials{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		CSRFKey:     csrfKey,
	}, nil
}

type credentialsData struct{}

func (r *credentialsData) Get(ctx context.Context) (*config.Credentials, error) {
	var e credentialsEntity
	k := credentialsKey(ctx, "DEFAULT")
	if err := datastore.Get(ctx, k, &e); err != nil {
		return nil, errors.Wrapf(err, "error while getting the entity")
	}
	csrfKey, err := base64.StdEncoding.DecodeString(e.CSRFKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error while decoding base64")
	}
	return &config.Credentials{
		GitHubToken: e.GitHubToken,
		CSRFKey:     csrfKey,
	}, nil
}

func credentialsKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, "Credentials", name, 0, nil)
}

type credentialsEntity struct {
	GitHubToken string // (required)
	CSRFKey     string // (required) base64 encoded string of 32 bytes key
}

# gradleupdate [![CircleCI](https://circleci.com/gh/int128/gradleupdate.svg?style=shield)](https://circleci.com/gh/int128/gradleupdate)

This service provides continuous update of Gradle Wrapper.
It automatically sends a pull request with the latest Gradle version.


## Getting Started

TODO


## Contributions

This is an open source software.
Feel free to open issues and pull requests.

### Architecture

This application is written in Go and designed for App Engine.
It consists of the following packages:

- `main` - Bootstraps the application and wires dependencies.
- `handlers` - Handles requests.
- `templates` - Renders pages.
- `usecases` - Provides application use cases.
- `domain` - Provides domain of weather forecast.
- `gateways` - Provides conversion between domain models and external models.
- `infrastructure` - Invokes external APIs.

### Development

Install dependencies:

```sh
brew install go
brew cask install google-cloud-sdk
gcloud components install app-engine-go

go get -u github.com/golang/mock/mockgen
go get -u github.com/valyala/quicktemplate/qtc
go get -u github.com/cortesi/modd/cmd/modd
```

Run the local server:

```sh
GITHUB_TOKEN="$GITHUB_TOKEN" CSRF_KEY="MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=" dev_appserver.py .
```

Run the mock server with `handlers` and `templates`:

```sh
modd
```

Regenerate templates and interface mocks:

```sh
go generate -v ./...
```

### Deployment

You need to create the following entity on Cloud Datastore:

- Kind = `Config`
- Key (string) = `DEFAULT`
- `GitHubToken` (string) = your GitHub token
- `CSRFKey` (string) = base64 encoded string of 32 bytes key

You can generate `CSRFKey` by the following command:

```sh
dd if=/dev/random bs=32 count=1 | base64
```

Deploy:

```sh
gcloud app deploy --project=gradleupdate
```

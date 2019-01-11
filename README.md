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

You need to create `.env.yaml` as follows:

```yaml
env_variables:
  GITHUB_TOKEN: YOUR_GITHUB_TOKEN
```

Run the local server:

```sh
dev_appserver.py .
```

Run the mock server (for development of templates):

```sh
modd
```

Deploy:

```sh
gcloud app deploy --project=gradleupdate
```

Regenerate templates and interface mocks:

```sh
go generate -v ./...
```

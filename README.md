# gradleupdate [![CircleCI](https://circleci.com/gh/int128/gradleupdate.svg?style=shield)](https://circleci.com/gh/int128/gradleupdate)

This provides continuous update of Gradle in your GitHub repositories.


## Getting Started

You need to add the following badge to README in a repository.

```markdown
[![Gradle Status](https://gradleupdate.appspot.com/YOUR/REPO/status.svg)](https://gradleupdate.appspot.com/YOUR/REPO/status)
```

And then, gradleupdate will send a pull request for the latest version of Gradle wrapper if it is out-of-dated.

You can turn off updates by removing the badge.


## How it works

gradleupdate checks the following preconditions for each repository:

- `README.md` should have the badge.
- `gradle/wrapper/gradle-wrapper.properties` should have a valid version string.

If the repository has out-of-dated Gradle version, gradleupdate will send a pull request.

If the pull request already exists, gradleupdate will just rebase the head branch.


## Contributions

This is an open source software.
Feel free to open issues and pull requests.

### Architecture

This application is written in Go and designed for App Engine.
It consists of the following packages:

- `main` - Bootstraps the application.
- `di` - Wires dependencies.
- `handlers` - Handles requests.
- `templates` - Renders pages.
- `usecases` - Provides application use cases.
- `domain` - Provides domain of Git, Gradle and gradleupdate.
- `gateways` - Provides conversion between domain models and external models.
- `infrastructure` - Invokes Gradle Services API, GitHub v3 API and GitHub v4 API.

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

You need to set up your Credentials.
See [gateways/credentials.go](gateways/credentials.go) for details.

You can set the feature toggles.
See [gateways/toggles.go](gateways/toggles.go) for details.

Deploy:

```sh
gcloud app deploy --project=gradleupdate
```

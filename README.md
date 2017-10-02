# GradleUpdate [![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/status.svg)](https://gradleupdate.appspot.com/int128/gradleupdate/status)

Automatic Gradle Updater.


## How to Run

Create `.env` file.

```properties
SYSTEM_GITHUB_ACCESS_TOKEN=
GITHUB_OAUTH_CLIENT_ID=
GITHUB_OAUTH_CLIENT_SECRET=
```

Google Cloud SDK is required.

```sh
brew cask install google-cloud-sdk
gcloud components install app-engine-java

# Run dev server
./gradlew appengineRun

# Deploy
./gradlew appengineDeploy
```


## How to Deploy

Configure following environment variables on CircleCI.

- `GCP_SERVICE_ACCOUNT_KEY` - Base64 encoded service account key
- `DOTENV` - Base64 encoded `.env` for production

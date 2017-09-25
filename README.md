# GradleUpdate [![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/status.svg)](https://gradleupdate.appspot.com/int128/gradleupdate/status)

Automatic Gradle Updater.


## How to Run

Create `.env` file.

```properties
GRADLEUPDATE_GITHUB_ACCESS_TOKEN=
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

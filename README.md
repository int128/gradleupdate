# GradleUpdate [![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/status.svg)](https://gradleupdate.appspot.com/int128/gradleupdate/status)

Automatic Gradle Updater.


## How to Run

Create `.properties` file into the project.

```properties
gradleupdate.github.accessToken=Personal access token for local development
security.oauth2.client.clientId=OAuth client ID for local development
security.oauth2.client.clientSecret=OAuth client secret for local development
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

```properties
GCP_SERVICE_ACCOUNT_KEY=Base64 encoded service account key
SYSTEM_GITHUB_ACCESS_TOKEN=Personal access token for production
GITHUB_OAUTH_CLIENT_ID=OAuth client ID for production
GITHUB_OAUTH_CLIENT_SECRET=OAuth client secret for production
```

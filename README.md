Gaelyk + React on App Engine [![CircleCI](https://circleci.com/gh/int128/gaelyk-react-starter.svg?style=svg)](https://circleci.com/gh/int128/gaelyk-react-starter) [![Gradle Status](https://gradleupdate.appspot.com/int128/gaelyk-react-starter/status.svg?branch=master)](https://gradleupdate.appspot.com/int128/gaelyk-react-starter/status)
============================

A template project with Gaelyk and React on App Engine.

* React
* Gaelyk
* App Engine
* Spock
* Webpack
* Babel
* Gradle

How to Run
----------

Build and run App Engine development server.

```bash
npm install
npm run build
./gradlew --daemon appengineRun
./gradlew appengineStop
```

We can run Webpack development server instead.

```bash
npm start
```

How to Deploy
-------------

Push the master branch and Circle CI will deploy the app to App Engine.

A service account key should be provided during CI.
Open [Google Cloud Platform Console](https://console.cloud.google.com/iam-admin/serviceaccounts) and create a service account.
Then, encode the JSON key as follows and store it into the environment variable `APPENGINE_KEY` on Circle CI.

```bash
base64 -b0 appengine-key.json
```

Build system
------------

Sources:

* Frontend
  * `/src/main/js` - JSX and Less code
  * `/static` - Static files
* Backend
  * `/src/main/groovy` - Production code
  * `/src/main/groovlet` - Production code (Groovlet)
  * `/src/main/config` - Configuration files
  * `/src/test/groovy` - Test code

Artifacts:

Destination                             | Source                                | Builder
----------------------------------------|---------------------------------------|---------
`/build/exploded-app`                   | `/static`                             | Webpack
`/build/exploded-app/react.min.js`      | dependencies                          | Webpack
`/build/exploded-app/app.js`            | `/src/main/js` and dependencies       | Webpack
`/build/exploded-app/WEB-INF`           | `/src/main/config`                    | Gradle
`/build/exploded-app/WEB-INF/lib`       | `/src/main/groovy` and dependencies   | Gradle
`/build/exploded-app/WEB-INF/groovy`    | `/src/main/groovlet`                  | Gradle and Webpack
`/build/local_db.bin`                   | -                                     | App Engine Dev Server

App Engine Groovy Template [![Build Status](https://travis-ci.org/int128/gradle-appengine-blank.svg)](https://travis-ci.org/int128/gradle-appengine-blank)
==========================

A template project with Groovy on Google App Engine.


* Groovy
* Gradle
  * App Engine support with [gradle-appengine-plugin](https://github.com/GoogleCloudPlatform/gradle-appengine-plugin)
  * Groovy lang support
  * IntelliJ IDEA support


Prepare
-------

Java 7 or later is required.

Install App Engine SDK.

```sh
gcloud components update gae-java
```

Set environment variables in `.bashrc` or `.zshrc`.

```sh
export APPENGINE_HOME=$HOME/Library/google-cloud-sdk/platform/appengine-java-sdk
```

Run
---

Run the development server.

```bash
./gradlew appengineRun

./gradlew appengineStop
```

Deploy
------

Deploy the application to the production platform.

```bash
./gradlew appengineUpdate
```

Structure
---------

* `src/main/groovy/` - Groovy sources of the product
* `src/test/groovy/` - Groovy sources of the test
* `static/public/` - Static files
* `static/target/` - Compiled assets
* `build.gradle` - Build settings

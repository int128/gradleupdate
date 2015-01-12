App Engine Hello World [![Build Status](https://travis-ci.org/int128/gradle-appengine-blank.svg)](https://travis-ci.org/int128/gradle-appengine-blank)
======================

An App Engine application with Gradle.


Architecture
------------

* Product
  * Groovy
  * Blank implementation of the router
* Build system
  * Gradle with App Engine plugin
  * Groovy lang support
  * IntelliJ IDEA support

How to use
----------

### Setting up environment

Java 7 or later is required.

Install App Engine SDK.

```bash
brew install app-engine-sdk-java
```

And set environment variables in `.bashrc` or `.zshrc`.

```bash
export APPENGINE_HOME=/usr/local/Cellar/app-engine-java-sdk/x.y.z/libexec
```

### Run app

Run the development server.

```bash
./gradlew appengineRun
```

Stop the development server.

```bash
./gradlew appengineStop
```

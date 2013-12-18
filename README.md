App Engine Blank Project
========================

Template project of App Engine application.


How to use
----------

`git clone` the repository and rename it.

Open `src/main/webapp/WEB-INF/appengine-web.xml` and change application id.

```xml
<appengine-web-app xmlns="http://appengine.google.com/ns/1.0">
    <application>myapp</application>
```

Then, invoke the gradle wrapper.

```bash
./gradlew appengineRun
```


Features
--------

This project contains these features:

  * Blank implementation of the router
  * Continuous integration support on Travis CI
  * Gradle Wrapper
  * `.gitignore` for Gradle, IDEA and Eclipse


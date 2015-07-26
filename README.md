# Gradle Update


## User Experience

A developer can keep the latest Gradle wrapper of repositories.

At first, you add a star on this repository to turn on update.
Then, we check version of Gradle wrapper in your repositories and open pull requests if needed.

When an new version of Gradle is released, we open pull requests for your repositories.


### Implementation

When an user added a star on this repository,

1. Receive an event of starred by _user_ from GitHub. Queue a request for _user_.
2. Get a list of repositories of _user_. Queue requests for each _repositories_.
3. Check version of Gradle wrapper in the _repository_. Queue a request for update if needed.
4. Invoke a worker for update of the _repository_.

When an new version of Gradle is released,

1. Periodically check version of the latest Gradle.
   If an new version is released, queue a request for users.
2. Get a list of users who starred on this repository. Queue requests for each _users_.


## Architecture

* Frontend
  * GitHub Pages (Jekyll)
  * HTML5 based Web Application
  * CoffeeScript
  * Vue.js
  * Bootswatch
* Backend
  * Google App Engine (JavaVM)
  * Gaelyk
  * Groovy 2.4





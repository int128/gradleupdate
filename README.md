# Gradle Update

Gradle Update will update Gradle wrapper of your repositories when the new Gradle is released.

## How to turn on update

Click [★ Star] button on [the repository page](https://github.com/int128/gradleupdate).



Gradle Update checks your repositories and sends pull requests with the latest Gradle wrapper if needed.
When an new version of Gradle is released, we will send pull requests with it.


## Architecture

* Backend
  * Google App Engine (JavaVM)
  * Gaelyk
  * Groovy 2.4


### Implementation

When an user added a star on this repository,

1. Receive an event of starred by _user_ from GitHub. Queue a request for _user_.
2. Get a list of repositories of _user_. Queue requests for each _repositories_.
3. Check version of Gradle wrapper in the _repository_. Queue a request for update if needed.
4. Fork and send a pull request for the _repository_.

When an new version of Gradle is released,

1. Periodically check version of the latest Gradle.
   If an new version is released, queue a request for users.
2. Get a list of stargazers of this repository. Queue requests for each _users_.

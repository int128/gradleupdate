# Gradle Update

Gradle Update keeps the latest build system. It updates Gradle Wrapper of your repositories by a pull request when a new version of Gradle is released.


## How to use

Star [the repository](https://github.com/int128/gradleupdate) to enable Gradle Update.

<img src="https://cloud.githubusercontent.com/assets/321266/9202088/176d83d6-408b-11e5-96dd-c138322fde60.png">

You will receive pull requests as follows.

<img src="https://cloud.githubusercontent.com/assets/321266/9202273/0e94da60-408c-11e5-83e9-594c9fbdcd42.png">

Open a pull request and check the content of new Gradle Wrapper.

<img src="https://cloud.githubusercontent.com/assets/321266/9202364/70fd5a6a-408c-11e5-9cc6-4a7a8f9ccfa8.png">

Merge the pull request if all tests have passed.

When a new version of Gradle is released, we will send pull requests for you.


## Architecture

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

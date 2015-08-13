# Gradle Update ![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/gradle.svg?branch=master)

Gradle Update automatically delivers the new Gradle Wrapper to your repositories by pull requests when a new version of Gradle is released.

The latest build system brings much benefit such as performance and bug fixes.


## Getting started

**Step 1:** Star [this repository](https://github.com/int128/gradleupdate) to enable Gradle Update.

<img src="https://cloud.githubusercontent.com/assets/321266/9202088/176d83d6-408b-11e5-96dd-c138322fde60.png">

**Step 2:** You will receive pull requests soon if you have any repositories which contain Gradle Wrapper.

<img src="https://cloud.githubusercontent.com/assets/321266/9202273/0e94da60-408c-11e5-83e9-594c9fbdcd42.png">

**Step 3:** Open a pull request and check the content of new Gradle Wrapper.

<img src="https://cloud.githubusercontent.com/assets/321266/9202364/70fd5a6a-408c-11e5-9cc6-4a7a8f9ccfa8.png">

Merge the pull request if all tests have passed.

When a new version of Gradle is released, Gradle Update will send pull requests for you.


## Origin

Gradle Update checks version of Gradle from [services.gradle.org](https://services.gradle.org).

Gradle Update creates a pull request for the latest version of Gradle Wrapper from [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper) continuously maintained by Gradle Update.


## Contribution

Gradle Update is an open source software licensed under the Apache License Version 2.0. Feel free to open issues or pull requests.


### Architecture

* Groovy 2.4
* Gaelyk
* Google App Engine (JavaVM)

Git operations are performed via GitHub API. It requires no filesystem or git command.

All operations are designed to be transactional and idempotence. Any exception such as HTTP error may occur during an operation but will be recovered by retry of the task queue.


### Design

When a user starred this repository,

1. Receive an event from GitHub. Queue a task for _user_.
2. In a task of _user_, get a list of repositories of _user_. Queue tasks for each _repositories_.
3. In a task of _repository_, check version of Gradle Wrapper of the _repository_. If it is out-of-dated, queue a task for update the _repository_.
4. In a task of updating _repository_, fork the _repository_ and send a pull request for the latest Gradle Wrapper from [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).

Periodically,

1. Check version of the latest Gradle.

When an new version of Gradle is released,

1. Get a list of stargazers of this repository. Queue tasks for each _user_.
2. Same as the starred event.

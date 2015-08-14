# Gradle Update [![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/status.svg?branch=master)](https://gradleupdate.appspot.com/int128/gradleupdate/status)

Gradle Update provides the latest Gradle Wrapper by pull requests for your repositories on GitHub.

The latest build system brings much benefit such as better performance and bug fixes.


## How to Use

Just 3 steps to turn on Gradle Update.

**Step 1:** Add a star to Gradle Update repository.

<img src="https://cloud.githubusercontent.com/assets/321266/9202088/176d83d6-408b-11e5-96dd-c138322fde60.png">

Gradle Update checks version of Gradle Wrapper in your public repositories. It may take a few minutes.

**Step 2:** Open [notifications page](https://github.com/notifications). You will receive pull requests for updating Gradle Wrapper.

<img src="https://cloud.githubusercontent.com/assets/321266/9202273/0e94da60-408c-11e5-83e9-594c9fbdcd42.png">

**Step 3:** Merge a pull request if all tests have passed.

<img src="https://cloud.githubusercontent.com/assets/321266/9202364/70fd5a6a-408c-11e5-9cc6-4a7a8f9ccfa8.png">

Once starred, you will receive pull requests when a new version of Gradle is released.


**Step 4:** Add a badge to indicate the build system is up-to-date.

[![Gradle Status](https://gradleupdate.appspot.com/int128/latest-gradle-wrapper/status.svg?branch=master)](https://gradleupdate.appspot.com/int128/latest-gradle-wrapper/status)

Add following line into `README.md` in your repository. Replace `USER` and `REPO` with proper one.

```markdown
[![Gradle Status](https://gradleupdate.appspot.com/USER/REPO/status.svg?branch=master)](https://gradleupdate.appspot.com/USER/REPO/status)
```


## How Gradle Update works

When a user starred this repository,

1. Receive an event from GitHub and queue a task for _user_.
2. In a task for _user_,
  1. Get a list of repositories of _user_.
  2. Queue tasks for each _repositories_.
3. In a task for _repository_,
  1. Check version of Gradle Wrapper in the _repository_.
  2. If it is out-of-date, queue a task for updating the _repository_.
  3. If it is up-to-date or has no Gradle Wrapper, do nothing.
4. In a task for updating _repository_,
  1. Fork the _repository_.
  2. Create a commit with the latest Gradle Wrapper from [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
  3. Create a pull request.

Periodically,

1. Check the latest version of Gradle from [services.gradle.org](https://services.gradle.org).

When a new version of Gradle is found,

1. Trigger updating Gradle Wrapper on [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
2. Wait until the new version is available on [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
3. Get a list of stargazers of this repository and queue tasks for each _stargazer_.


### Architecture

* Groovy 2.4
* Gaelyk
* Google App Engine (JavaVM)

Git operations are performed via GitHub API. It requires no filesystem or git command.

All operations are designated to be transactional and idempotence. Any exception such as HTTP error may occur during an operation but will be recovered by retry of the task queue.


## Contribution

Gradle Update is an open source software licensed under the Apache License Version 2.0. Feel free to open issues or pull requests.

# Gradle Update [![Gradle Status](https://gradleupdate.appspot.com/int128/gradleupdate/status.svg?branch=master)](https://gradleupdate.appspot.com/int128/gradleupdate/status)

Gradle Update keeps the latest Gradle Wrapper by pull requests for your repositories on GitHub.

The latest build system brings much benefit such as better performance and bug fixes.


## How to Use

Just 3 steps to turn on Gradle Update.

----

**Step 1:** Add a star to Gradle Update repository.

<img src="https://cloud.githubusercontent.com/assets/321266/9202088/176d83d6-408b-11e5-96dd-c138322fde60.png">

Gradle Update receives an event from GitHub and checks version of Gradle Wrapper in your repositories. It can access to only public repositories.

----

**Step 2:** Open [notifications page](https://github.com/notifications). You will receive pull requests for updating Gradle Wrapper after a few minutes.

<img src="https://cloud.githubusercontent.com/assets/321266/9202273/0e94da60-408c-11e5-83e9-594c9fbdcd42.png">

----

**Step 3:** Merge a pull request if all tests have passed.

<img src="https://cloud.githubusercontent.com/assets/321266/9202364/70fd5a6a-408c-11e5-9cc6-4a7a8f9ccfa8.png">

Once starred, you will receive pull requests when a new version of Gradle is released. You can turn off Gradle Update by removing the star at any time.

----

**Optional:** Add a badge to indicate the build system is up-to-date, e.g.:

[![Gradle Status](https://gradleupdate.appspot.com/int128/latest-gradle-wrapper/status.svg?branch=master)](https://gradleupdate.appspot.com/int128/latest-gradle-wrapper/status)

Add the following line to `README.md` in your repositories. Replace `USER` with the user name, and `REPO` with the repository name.

```markdown
[![Gradle Status](https://gradleupdate.appspot.com/USER/REPO/status.svg?branch=master)](https://gradleupdate.appspot.com/USER/REPO/status)
```


## How Gradle Update works

When a _user_ added a star to Gradle Update repository,

1. Receive an event from GitHub and queue a task for _user_.
2. In a task for the _user_,
  1. Get a list of repositories of _user_.
  2. Queue tasks for each _repositories_.
3. In a task for the _repository_,
  1. Check version of Gradle Wrapper in the _repository_.
  2. If following conditions are satisfied, queue a task for updating the _repository_.
    1. The _repository_ contains Gradle Wrapper.
    2. Gradle Wrapper of the _repository_ is out-of-date.
    3. There is no pull request already sent.
4. In a task for updating the _repository_,
  1. Fork the _repository_.
  2. Create a commit with the latest Gradle Wrapper from [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
  3. Create a pull request.

Periodically,

1. Check the latest version of Gradle from [services.gradle.org](https://services.gradle.org).
2. If a new version of Gradle is found,
  1. Trigger updating Gradle Wrapper on [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
  2. Wait until the new version is available on [latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
  3. Get a list of stargazers of this repository and queue tasks for each _user_.

When a badge for the _repository_ is requested,

1. Check the version of Gradle Wrapper in the _repository_.
2. If it is up-to-date, respond the green image.
3. If it is out-of-date, respond the red image.
4. If it has no Gradle Wrapper, respond the grey image.


### Architecture

* Groovy 2.4
* [Gaelyk](gaelyk.appspot.com)
* [groovy-wslite](https://github.com/jwagenleitner/groovy-wslite)
* Google App Engine (JavaVM)

Git operations are performed via GitHub API. It requires no filesystem or git command.

All operations are performed on Task Queue and designed to be transactional and idempotence. Any exception such as HTTP error may occur during an operation but will be recovered by retrying.


## Contribution

Gradle Update is an open source software licensed under the Apache License Version 2.0. Feel free to open issues or pull requests.

import infrastructure.GitHub

final from = params.from
assert from instanceof String

final fullName = params.full_name
assert fullName instanceof String

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()

final title = "Gradle $gradleVersion"
final body = """
[Gradle $gradleVersion](https://gradle.org/docs/$gradleVersion/release-notes) is available now.

This pull request updates Gradle wrapper and build.gradle in the repository.
Please merge this if all tests are passed with the latest Gradle.
"""

log.info("Creating a pull request from $from into $fullName")
final pullRequest = gitHub.createPullRequest(fullName, 'master', from, title, body)
assert pullRequest

log.info("Pull request #${pullRequest.number} has been created on ${pullRequest.html_url}")

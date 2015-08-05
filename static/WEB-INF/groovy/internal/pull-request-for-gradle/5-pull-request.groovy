import infrastructure.GitHub

final fromUser = params.from_user
final fromBranch = params.from_branch
final intoRepo = params.into_repo
final intoBranch = params.into_branch
final gradleVersion = params.gradle_version
assert fromUser instanceof String
assert fromBranch instanceof String
assert intoRepo instanceof String
assert intoBranch instanceof String
assert gradleVersion instanceof String

final gitHub = new GitHub()

final title = "Gradle $gradleVersion"
final body = """
[Gradle $gradleVersion](https://gradle.org/docs/$gradleVersion/release-notes) is available now.

This pull request updates Gradle wrapper and build.gradle in the repository.
Please merge this if all tests are passed with the latest Gradle.
"""

log.info("Creating a pull request from $fromBranch into $intoRepo:$intoBranch")
final pullRequest = gitHub.createPullRequest(intoRepo, intoBranch, "$fromUser:$fromBranch", title, body)
assert pullRequest

log.info("Pull request #${pullRequest.number} has been created on ${pullRequest.html_url}")

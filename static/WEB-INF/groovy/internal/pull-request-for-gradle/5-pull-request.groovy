import gradle.Repository
import gradle.Stargazers

import static util.RequestUtil.relativePath

assert params.from_user
assert params.from_branch
assert params.into_repo
assert params.into_branch
assert params.gradle_version

final stargazers = new Stargazers()

final title = "Gradle $params.gradle_version"
final body = """
[Gradle $params.gradle_version](https://gradle.org/docs/$params.gradle_version/release-notes) is available now.

This pull request updates Gradle wrapper and build.gradle in the repository.
Merge it if all tests passed with the latest Gradle.

Automatic pull request can be turned off by unstar [gradleupdate repository](${stargazers.htmlUrl}).
"""

final pullRequest = new Repository(params.into_repo).createPullRequest(params.into_branch, params.from_user, params.from_branch, title, body)

log.info("Pull request #${pullRequest.number} has been created on ${pullRequest.html_url}")
defaultQueue.add(
        url: relativePath(request, '6-done.groovy'),
        params: [
                gradleVersion: params.gradle_version,
                url: pullRequest.url,
                htmlUrl: pullRequest.html_url,
                createdAt: pullRequest.created_at,
                repo: pullRequest.base.repo.name,
                owner: pullRequest.base.repo.owner.login,
                ownerId: pullRequest.base.repo.owner.id,
        ])

import gradle.Repository
import gradle.TemplateRepository

import static util.RequestUtil.relativePath

assert params.from_user
assert params.from_repo
assert params.from_branch
assert params.into_repo
assert params.into_branch
assert params.gradle_version

final templateRepository = new TemplateRepository()
final fromRepository = new Repository(params.from_repo)

log.info("Creating a tree with the latest Gradle wrapper on $params.from_repo")
final treeForGradleWrapper = fromRepository.createTreeForGradleWrapper(templateRepository)

log.info("Creating a tree with build.gradle for $params.gradle_version on $params.from_repo")
final treeForBuildGradle = fromRepository.createTreeForBuildGradle(params.into_branch, params.gradle_version)

final tree = treeForGradleWrapper + treeForBuildGradle

final created = fromRepository.createBranch(params.from_branch, params.into_branch, "Gradle $params.gradle_version", tree)
if (!created) {
    log.info("Already up-to-date version $params.gradle_version, skip")
}

log.info("Queue sending a pull request into $params.into_repo:$params.into_branch from $params.from_user:$params.from_branch")
defaultQueue.add(
        url: relativePath(request, '5-pull-request.groovy'),
        params: [
                from_user: params.from_user,
                from_branch: params.from_branch,
                into_repo: params.into_repo,
                into_branch: params.into_branch,
                gradle_version: params.gradle_version,
        ],
        countdownMillis: 1000)

import gradle.Repository

assert params.full_name

final repository = new Repository(params.full_name as String)
if (repository.queryIfHasGradleWrapper()) {
    log.info("Repository ${params.full_name} has Gradle wrapper, sending a pull request for the latest one")
    repository.sendPullRequestForLatestGradle()
} else {
    log.info("Repository ${params.full_name} does not have Gradle wrapper, skip")
}

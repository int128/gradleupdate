import infrastructure.GitHub

final gradleVersion = params.gradle_version
assert gradleVersion

log.warning("Closing all pull requests for Gradle $gradleVersion")

final gitHub = new GitHub()

log.info("Fetching our repositories")
final repositories = gitHub.fetchRepositories('gradleupdate')
assert repositories instanceof List

repositories.findAll { repository ->
    assert repository.fork instanceof Boolean
    repository.fork
}.each { repository ->
    final fullName = repository.full_name
    assert fullName instanceof String

    log.info("Removing branch gradle-$gradleVersion of $fullName")
    gitHub.removeBranch(fullName, "gradle-$gradleVersion")
}

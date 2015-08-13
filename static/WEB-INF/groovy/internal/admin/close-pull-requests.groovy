import gradle.Repository

final gradleVersion = params.gradle_version
assert gradleVersion

log.warning("Closing all pull requests for Gradle $gradleVersion")

Repository.fetchRepositories('gradleupdate').findAll { repository ->
    assert repository.fork instanceof Boolean
    repository.fork
}.each { repository ->
    final fullName = repository.full_name
    assert fullName instanceof String
    new Repository(fullName).removeBranch("gradle-$gradleVersion")
}

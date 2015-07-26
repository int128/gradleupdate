assert params.gradleVersion

final repositories = null //TODO

repositories.each { repo ->
    log.info("Queue updating the repository: $repo.fullName")
    defaultQueue.add(
            url: '/internal/update-gradle-of-user-repository.groovy',
            params: [fullName: repo.fullName, gradleVersion: params.gradleVersion])
}

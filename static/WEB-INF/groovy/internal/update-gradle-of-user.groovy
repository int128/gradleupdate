import infrastructure.GitHub

assert params.user

final gitHub = new GitHub()
final repositories = gitHub.getRepositories(params.user)

repositories.each { repo ->
    log.info("Queue updating the repository: $repo.full_name")
    defaultQueue.add(
            url: '/internal/update-gradle-of-repository.groovy',
            params: [full_name: repo.full_name])
}

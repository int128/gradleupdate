import infrastructure.GitHub

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final stargazer = params.stargazer
assert stargazer instanceof String

final gitHub = new GitHub()

log.info("Fetching repositories of stargazer $stargazer")
final repositories = gitHub.fetchRepositories(stargazer)

repositories.each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/',
            params: [full_name: repo.full_name, gradle_version: gradleVersion])
}

import infrastructure.GitHub
import service.GradleVersionService

final stargazer = params.stargazer
assert stargazer instanceof String

log.info("Fetching version of the latest Gradle")
final gradleVersion = new GradleVersionService().queryStableVersion()

log.info("Fetching repositories of stargazer $stargazer")
final gitHub = new GitHub()
final repositories = gitHub.getRepositories(stargazer)

repositories.each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/',
            params: [full_name: repo.full_name, gradle_version: gradleVersion])
}
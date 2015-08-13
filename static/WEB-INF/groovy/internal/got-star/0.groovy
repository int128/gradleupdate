import gradle.VersionWatcher
import infrastructure.GitHub

final stargazer = params.stargazer
assert stargazer instanceof String

log.info("Fetching version of the latest Gradle")
final gradleVersion = new VersionWatcher().fetchStableVersion()

log.info("Fetching repositories of stargazer $stargazer")
final gitHub = new GitHub()
final repositories = gitHub.fetchRepositories(stargazer)

repositories.each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/0.groovy',
            params: [
                    full_name: repo.full_name,
                    branch: repo.default_branch,
                    gradle_version: gradleVersion,
            ])
}

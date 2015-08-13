import gradle.Repository
import gradle.VersionWatcher

final stargazer = params.stargazer
assert stargazer instanceof String

log.info("Fetching version of the latest Gradle")
final gradleVersion = new VersionWatcher().fetchStableVersion()

Repository.fetchRepositories(stargazer).each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/0.groovy',
            params: [
                    full_name: repo.full_name,
                    branch: repo.default_branch,
                    gradle_version: gradleVersion,
            ])
}

import gradle.Repository
import infrastructure.GitHub
import service.GradleVersionService

final fullName = params.full_name
assert fullName instanceof String

final gitHub = new GitHub()
final repository = new Repository(fullName, gitHub)

log.info("Checking Gradle version of repository $fullName")
final gradleVersion = repository.queryGradleWrapperVersion()
if (gradleVersion == null) {
    log.info("$fullName does not have Gradle wrapper")
    return
}

final latestGradleVersion = new GradleVersionService().queryStableVersion()
if (gradleVersion == latestGradleVersion) {
    log.info("$fullName has the latest Gradle wrapper $gradleVersion")
    return
}

log.info("$fullName has obsolete Gradle wrapper $gradleVersion " +
        "while latest is $latestGradleVersion, queue updating")
defaultQueue.add(
        url: '/internal/pull-request-for-gradle/',
        params: [full_name: fullName])

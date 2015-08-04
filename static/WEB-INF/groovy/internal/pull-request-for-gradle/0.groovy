import gradle.Repository
import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()
final repository = new Repository(fullName, gitHub)

log.info("Checking Gradle version of repository $fullName")
final gradleWrapperVersion = repository.queryGradleWrapperVersion()

if (gradleWrapperVersion == null) {
    log.info("$fullName does not have Gradle wrapper")
    return
}
if (gradleWrapperVersion == gradleVersion) {
    log.info("$fullName has the latest Gradle wrapper $gradleVersion")
    return
}

log.info("$fullName has obsolete Gradle wrapper $gradleWrapperVersion " +
        "while latest is $gradleVersion, so queue updating")
defaultQueue.add(
        url: relativePath(request, '1-fork.groovy'),
        params: [full_name: fullName, gradle_version: gradleVersion])

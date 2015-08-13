import gradle.Repository
import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
final branch = params.branch
final gradleVersion = params.gradle_version
assert fullName instanceof String
assert branch instanceof String
assert gradleVersion instanceof String

final gitHub = new GitHub()
final repository = new Repository(fullName, gitHub)

log.info("Checking Gradle version of repository $fullName")
final gradleWrapperVersion = repository.fetchGradleWrapperVersion(branch)

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
        params: [
                into_repo: fullName,
                into_branch: branch,
                gradle_version: gradleVersion
        ])

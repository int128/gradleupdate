import gradle.Repository

import static util.RequestUtil.relativePath

assert params.full_name
assert params.branch
assert params.gradle_version

final repository = new Repository(params.full_name)
final gradleWrapperVersion = repository.fetchGradleWrapperVersion(params.branch)

if (gradleWrapperVersion == null) {
    log.info("$params.full_name does not have Gradle wrapper")
    return
}
if (gradleWrapperVersion == params.gradle_version) {
    log.info("$params.full_name has the latest Gradle wrapper $params.gradle_version")
    return
}

log.info("$params.full_name has obsolete Gradle wrapper $gradleWrapperVersion " +
        "while latest is $params.gradle_version, so queue updating")
defaultQueue.add(
        url: relativePath(request, '1-fork.groovy'),
        params: [
                into_repo: params.full_name,
                into_branch: params.branch,
                gradle_version: params.gradle_version
        ])

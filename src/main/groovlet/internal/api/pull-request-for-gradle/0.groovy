import gradle.Repository

assert params.full_name
assert params.branch
assert params.gradle_version

final origin = new Repository(params.full_name)
final originGradleWrapperVersion = origin.fetchGradleWrapperVersion(params.branch)

if (originGradleWrapperVersion == null) {
    log.info("Repository $params.full_name does not have Gradle wrapper")
    return
}
if (originGradleWrapperVersion == params.gradle_version) {
    log.info("Repository $params.full_name already has the latest Gradle wrapper $params.gradle_version")
    return
}
log.info("Repository $params.full_name has obsolete Gradle wrapper $originGradleWrapperVersion")

final fork = origin.fork()

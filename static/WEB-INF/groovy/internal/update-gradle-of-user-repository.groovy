import infrastructure.GradleUpdateWorker

assert params.fullName
assert params.gradleVersion

final worker = new GradleUpdateWorker()

if (worker.queryGradleWrapperVersion() == params.gradleVersion) {
    worker.bumpUserRepository(params.fullName)
} else {
    response.sendError(404, "Gradle $params.gradleVersion does not found")
}

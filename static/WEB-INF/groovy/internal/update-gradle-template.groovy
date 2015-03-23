import infrastructure.GradleUpdateWorker

assert params.gradleVersion

final worker = new GradleUpdateWorker()
worker.bumpTemplate(params.gradleVersion)

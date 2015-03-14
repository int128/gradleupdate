import service.GradleUpdateWorker

assert params.'gradle-version'

final worker = new GradleUpdateWorker()
worker.bumpTemplate(params.'gradle-version')

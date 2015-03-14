import service.GradleUpdateWorker

assert params.repo

final worker = new GradleUpdateWorker()
worker.bumpUserRepository(params.repo)

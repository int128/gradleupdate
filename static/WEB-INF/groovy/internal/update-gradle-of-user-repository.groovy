import infrastructure.GradleUpdateWorker

assert params.repo

final worker = new GradleUpdateWorker()
worker.bumpUserRepository(params.repo)

import service.GradleUpdateWorker

assert params.full_name

final worker = new GradleUpdateWorker()
worker.bumpUserRepository(params.full_name)

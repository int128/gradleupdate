import infrastructure.GradleUpdateWorker

final worker = new GradleUpdateWorker()

response.contentType = 'text/plain'
println worker.queryGradleWrapperVersion()

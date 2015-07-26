import infrastructure.GradleUpdateWorker

final worker = new GradleUpdateWorker()
final stargazers = worker.queryStargazers()

stargazers.each { stargazer ->
    log.info("Queue updating repositories of user: ${stargazer.login}")
    defaultQueue.add(
            url: '/internal/update-gradle-of-user.groovy',
            params: [user: stargazer.login])
}

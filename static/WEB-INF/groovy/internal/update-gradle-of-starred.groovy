import service.Stargazers

final stargazers = new Stargazers().fetch()

stargazers.each { stargazer ->
    log.info("Queue updating repositories of user: ${stargazer.login}")
    defaultQueue.add(
            url: '/internal/update-gradle-of-user.groovy',
            params: [user: stargazer.login])
}

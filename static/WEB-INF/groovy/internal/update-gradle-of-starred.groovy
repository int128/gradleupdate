import infrastructure.GitHub
import service.Stargazers

final gitHub = new GitHub()
final stargazers = new Stargazers(gitHub).fetch()

stargazers.each { stargazer ->
    log.info("Queue updating repositories of user: ${stargazer.login}")
    defaultQueue.add(
            url: '/internal/update-gradle-of-user.groovy',
            params: [user: stargazer.login])
}

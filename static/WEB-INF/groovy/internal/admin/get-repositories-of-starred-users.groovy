import gradle.Stargazers
import infrastructure.GitHub

final gitHub = new GitHub()
final stargazers = new Stargazers(gitHub).fetch()

final repositories = stargazers.collect { stargazer ->
    [
        name: stargazer.login,
        repos: gitHub.fetchRepositories(stargazer.login)*.full_name
    ]
}

response.contentType = 'text/plain'
println repositories

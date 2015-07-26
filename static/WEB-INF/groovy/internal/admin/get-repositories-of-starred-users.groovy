import infrastructure.GitHub
import service.Stargazers

final gitHub = new GitHub()
final stargazers = new Stargazers().fetch()
final repositories = stargazers.collect { stargazer ->
    [
        name: stargazer.login,
        repos: gitHub.getRepositories(stargazer.login)*.full_name
    ]
}

response.contentType = 'text/plain'
println repositories

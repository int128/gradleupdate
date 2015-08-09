import gradle.Stargazers
import infrastructure.GitHub

final gitHub = new GitHub()
final stargazers = new Stargazers(gitHub).fetch()

html.html {
    body {
        h1('Stargazers')
        ul {
            stargazers.each { stargazer ->
                li {
                    final repositories = gitHub.fetchRepositories(stargazer.login)
                    assert repositories instanceof List
                    a(href: stargazer.html_url, "$stargazer.login (${repositories.size()})")
                    ul {
                        repositories.each { repository ->
                            li {
                                a(href: repository.html_url, repository.full_name)
                            }
                        }
                    }
                }
            }
        }
    }
}

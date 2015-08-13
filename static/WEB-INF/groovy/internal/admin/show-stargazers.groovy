import gradle.Repository
import gradle.Stargazers

final stargazers = new Stargazers().fetch()

html.html {
    body {
        h1('Stargazers')
        ul {
            stargazers.each { stargazer ->
                li {
                    final repositories = Repository.fetchRepositories(stargazer.login)
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

import gradle.Repositories
import gradle.Stargazers

final stargazers = new Stargazers().fetchFirst().current

html.html {
    body {
        h1('Stargazers')
        ul {
            stargazers.each { stargazer ->
                li {
                    final repositories = new Repositories(stargazer.login as String).fetchFirst().current
                    a(href: stargazer.html_url, "$stargazer.login ($stargazer.id) (${repositories.size()})")
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

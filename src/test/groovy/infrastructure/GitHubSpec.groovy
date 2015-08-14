package infrastructure

import spock.lang.Specification

class GitHubSpec extends Specification {

    def "fetch stargazers should return a paged result"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def first = gitHub.fetchStargazersOfFirstPage('octocat/Spoon-Knife')

        then:
        first.current.size() > 0
        first.rel.next instanceof String
        first.rel.next.startsWith('https://api.github.com')

        when:
        def second = gitHub.fetchNextPage(first.rel.next)

        then:
        second.current.size() > 0
        second.rel.next instanceof String
        second.rel.next.startsWith('https://api.github.com')

        when:
        def third = gitHub.fetchNextPage(first.rel.next)

        then:
        third.current.size() > 0
        third.rel.next instanceof String
        third.rel.next.startsWith('https://api.github.com')
    }

    def "fetch stargazers for non-existent user should return null"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def page = gitHub.fetchStargazersOfFirstPage('0zCDPKcfjl1BWA6J/0zCDPKcfjl1BWA6J')

        then:
        page == null
    }

    def "fetch repositories should return a paged result"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def page = gitHub.fetchRepositoriesOfFirstPage('octocat')

        then:
        page.rel.next == null
        page.current instanceof List
        page.current.size() > 0
    }

    def "fetch repositories for non-existent user should return null"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def page = gitHub.fetchRepositoriesOfFirstPage('0zCDPKcfjl1BWA6J')

        then:
        page == null
    }

    def "fetchReference() should return refs and sha"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def ref = gitHub.fetchReference('octocat/Spoon-Knife', 'master')

        then:
        ref.ref == 'refs/heads/master'
        ref.object.type == 'commit'
        ref.object.sha =~ /[0-9a-z]{32}/
    }

}

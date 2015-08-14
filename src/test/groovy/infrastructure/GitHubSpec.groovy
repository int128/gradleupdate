package infrastructure

import spock.lang.Specification

class GitHubSpec extends Specification {

    def "fetchRepositories() should return a list"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def repos = gitHub.fetchRepositories('octocat')

        then:
        repos instanceof List
        repos.size() > 0
    }

    def "fetchRepositories(non-existent user) should return null"() {
        given:
        def gitHub = new GitHub(null)

        when:
        def repos = gitHub.fetchRepositories('0zCDPKcfjl1BWA6J')

        then:
        repos == null
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

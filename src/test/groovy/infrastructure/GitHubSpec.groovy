package infrastructure

import spock.lang.Specification

class GitHubSpec extends Specification {

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

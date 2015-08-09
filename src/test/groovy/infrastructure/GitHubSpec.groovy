package infrastructure

import spock.lang.Specification

class GitHubSpec extends Specification {

    def "getReference(master) should return refs and sha"() {
        given:
        def repository = new GitHub()

        when:
        def ref = repository.fetchReference('gradleupdate/Spoon-Knife', 'master')

        then:
        ref.ref == 'refs/heads/master'
        ref.object.type == 'commit'
        ref.object.sha =~ /[0-9a-z]{32}/
    }

    def "branch should be created and removed on each method"() {
        given:
        def repository = new GitHub()
        def branchName = "test${Math.random() * 1000 as int}"

        when:
        def ref = repository.createBranch('gradleupdate/Spoon-Knife', branchName, 'master')

        then:
        ref.ref == "refs/heads/$branchName"
        ref.object.type == 'commit'
        ref.object.sha =~ /[0-9a-z]{32}/

        when:
        def removed = repository.removeBranch('gradleupdate/Spoon-Knife', branchName)

        then:
        removed
    }

}

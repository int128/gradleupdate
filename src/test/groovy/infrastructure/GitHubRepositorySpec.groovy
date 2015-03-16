package infrastructure

import config.Credential
import spock.lang.Specification

class GitHubRepositorySpec extends Specification {

    def "getReference(master) should return refs and sha"() {
        given:
        def repository = new GitHubRepository('gradleupdate/Spoon-Knife', Credential.github)

        when:
        def ref = repository.getReference('master')

        then:
        ref.ref == 'refs/heads/master'
        ref.object.type == 'commit'
        ref.object.sha =~ /[0-9a-z]{32}/
    }

    def "branch should be created and removed on each method"() {
        given:
        def repository = new GitHubRepository('gradleupdate/Spoon-Knife', Credential.github)
        def branchName = "test${Math.random() * 1000 as int}"

        when:
        def ref = repository.createBranch(branchName, 'master')

        then:
        ref.ref == "refs/heads/$branchName"
        ref.object.type == 'commit'
        ref.object.sha =~ /[0-9a-z]{32}/

        when:
        def removed = repository.removeBranch(branchName)

        then:
        removed
    }

}

package infrastructure

import spock.lang.Specification

class GitHubWebhookSpec extends Specification {

    def 'A signature of GitHub webhook should be HMAC-SHA1'() {
        given:
        def secret = 'tCMh4MBuwIozkNyAY2yQOkeI4OIOc80wLodmi633puai8j4mvhdiT95HV0CrioXi'
        def signature = 'sha1=25c0e47af90de2194ec68c64cbaf6b1bab98eef1'
        def payload = GitHubWebhookSpec.class.getResourceAsStream('payload.json').bytes

        when:
        def valid = new GitHubWebhook(secret.bytes).validate(signature, payload)

        then:
        valid
    }

}

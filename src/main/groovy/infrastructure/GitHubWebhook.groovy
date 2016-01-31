package infrastructure

import model.Credential

import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec

import static model.Credential.CredentialKey.GitHubWebHookSecret

class GitHubWebhook {

    private final byte[] secret

    def GitHubWebhook() {
        secret = Credential.get(GitHubWebHookSecret).secret?.bytes
    }

    def GitHubWebhook(byte[] secret) {
        this.secret = secret
    }

    boolean validate(String signature, byte[] payload) {
        assert signature
        assert signature.startsWith('sha1=')
        def actual = signature.substring('sha1='.length()).decodeHex()

        assert payload
        def mac = Mac.getInstance('HmacSHA1')
        mac.init(new SecretKeySpec(secret, 'HmacSHA1'))
        def expected = mac.doFinal(payload)

        // constant time comparison for timing attack
        assert actual.length == expected.length
        (0..(actual.length - 1)).every { int i ->
            actual[i] == expected[i]
        }
    }

}

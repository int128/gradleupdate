package infrastructure

import model.Credential

import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec

class GitHubWebhook {

    private final byte[] secret

    def GitHubWebhook() {
        secret = Credential.getOrCreate('github-webhook').secret?.bytes
    }

    def GitHubWebhook(byte[] secret) {
        this.secret = secret
    }

    boolean validate(String signature, byte[] payload) {
        assert signature
        assert payload
        def mac = Mac.getInstance('HmacSHA1')
        mac.init(new SecretKeySpec(secret, 'HmacSHA1'))
        def expected = mac.doFinal(payload)
        signature == "sha1=${expected.encodeHex()}".toString()
    }

}

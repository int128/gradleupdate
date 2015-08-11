package infrastructure

import groovyx.net.http.HttpURLClient
import model.Credential

class GitHubOAuth {

    private final Credential credential

    def GitHubOAuth() {
        credential = Credential.getOrCreate('github')
    }

    def exchangeCodeAndToken(String code) {
        final client = new HttpURLClient(url: 'https://github.com/login/oauth/access_token')
        client.request(query: [
                client_id: credential.clientId,
                client_secret: credential.clientSecret,
                code: code
        ]).data
    }

}

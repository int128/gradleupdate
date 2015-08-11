package infrastructure

import groovyx.net.http.HttpURLClient
import model.Credential

class GitHubOAuth {

    final client = new HttpURLClient(url: 'https://github.com/login/oauth/access_token')

    def exchangeCodeAndToken(String code) {
        def clientId = Credential.getOrCreate('github-client-id')
        def clientKey = Credential.getOrCreate('github-client-key')
        client.request(query: [
                client_id: clientId.secret,
                client_secret: clientKey.secret,
                code: code
        ]).data
    }

}

package infrastructure

import model.Credential
import wslite.rest.RESTClient

class GitHubOAuth {

    final client

    def GitHubOAuth() {
        client = new RESTClient('https://github.com/login/oauth/access_token')
        client.httpClient.defaultHeaders += [Accept: 'application/json']
    }

    def exchangeCodeAndToken(String code) {
        def clientId = Credential.getOrCreate('github-client-id')
        def clientKey = Credential.getOrCreate('github-client-key')
        client.get(query: [
                client_id: clientId.secret,
                client_secret: clientKey.secret,
                code: code
        ]).json
    }

}

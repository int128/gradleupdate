package infrastructure

import entity.Credential
import wslite.rest.RESTClient

import static entity.Credential.CredentialKey.GitHubClientId
import static entity.Credential.CredentialKey.GitHubClientKey

class GitHubOAuth {

    final client = new RESTClient('https://github.com/login/oauth/access_token')

    def GitHubOAuth() {
        client.httpClient.defaultHeaders += [Accept: 'application/json']
    }

    def exchangeCodeAndToken(String code) {
        def clientId = Credential.get(GitHubClientId)
        def clientKey = Credential.get(GitHubClientKey)
        client.get(query: [
                client_id: clientId.secret,
                client_secret: clientKey.secret,
                code: code
        ]).json
    }

}

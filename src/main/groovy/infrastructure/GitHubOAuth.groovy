package infrastructure

import model.Credential
import wslite.rest.RESTClient

import static model.Credential.CredentialKey.GitHubClientId
import static model.Credential.CredentialKey.GitHubClientKey

class GitHubOAuth {

    final client

    def GitHubOAuth() {
        client = new RESTClient('https://github.com/login/oauth/access_token')
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

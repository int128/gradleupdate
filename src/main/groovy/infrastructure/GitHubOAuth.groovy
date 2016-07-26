package infrastructure

import entity.Credential
import wslite.rest.RESTClient

import static entity.Credential.CredentialKey.GitHubClientId
import static entity.Credential.CredentialKey.GitHubClientKey

class GitHubOAuth {

    static exchangeCodeAndToken(String code, String state, String redirectURI) {
        new RESTClient('https://github.com/login/oauth/access_token').get(query: [
                client_id: Credential.get(GitHubClientId).secret,
                client_secret: Credential.get(GitHubClientKey).secret,
                code: code,
                state: state,
                redirect_uri: redirectURI,
        ], headers: [
                Accept: 'application/json'
        ])
    }

}

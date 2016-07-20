package domain

import entity.Credential
import infrastructure.MemcacheHTTPClient
import infrastructure.RestAPI
import wslite.rest.RESTClient

import static entity.Credential.CredentialKey.GitHubToken

class GHSession {

    final RESTClient client = new RESTClient('https://api.github.com', new MemcacheHTTPClient())

    private def GHSession() {
        client.httpClient.defaultHeaders += ['User-Agent': 'gradleupdate']
    }

    private def GHSession(String oauthToken) {
        this()
        assert oauthToken
        client.httpClient.defaultHeaders += [Authorization: "token $oauthToken"]
    }

    static GHSession defaultToken() {
        new GHSession(Credential.get(GitHubToken).secret)
    }

    static GHSession byToken(String oauthToken) {
        new GHSession(oauthToken)
    }

    static GHSession noToken() {
        new GHSession()
    }

    @Lazy
    def repositories = { RestAPI.of(GHRepository, '/repos', this, client) }()

    GHRepository getRepository(String fullName) {
        repositories.get(fullName)
    }

}

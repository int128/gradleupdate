package infrastructure

import config.Credential
import groovyx.net.http.HttpURLClient

class GitHubOAuth {

    final client = new HttpURLClient(url: 'https://github.com/login/oauth/access_token')

    def exchange(String code) {
        client.request(query: [
            client_id: Credential.githubClientId,
            client_secret: Credential.githubClientSecret,
            code: code
        ]).data
    }

}

package infrastructure

import groovyx.net.http.ContentType
import groovyx.net.http.HttpResponseException
import groovyx.net.http.HttpURLClient

import static groovyx.net.http.Method.DELETE
import static groovyx.net.http.Method.POST

class GitHub {

    private final HttpURLClient client

    def GitHub(Map headers = [:]) {
        client = new HttpURLClient(url: 'https://api.github.com', headers: [
                'Authorization': "token $Credential.githubToken",
                'User-Agent': 'gradleupdate'
        ] + headers)
    }

    def whoAmI() {
        client.request(path: '/user').data
    }

    def getRepository(String repo) {
        client.request(path: "/repos/$repo").data
    }

    def createBranch(String repo, String branchName, String from) {
        def sha = getReference(repo, from).object.sha
        assert sha
        createReference(repo, branchName, sha)
    }

    boolean removeBranch(String repo, String branchName) {
        try {
            client.request(path: "/repos/$repo/git/refs/heads/$branchName", method: DELETE).success
        } catch (HttpResponseException e) {
            if (e.response.status == 422) {
                // API returns 422 if branch does not exist
                false
            } else {
                throw e
            }
        } catch (NullPointerException ignore) {
            // 204 No Content causes NPE due to the bug of HttpURLClient
            true
        }
    }

    def getReference(String repo, String branchName) {
        client.request(path: "/repos/$repo/git/refs/heads/$branchName").data
    }

    def createReference(String repo, String branchName, String shaRef) {
        client.request(path: "/repos/$repo/git/refs", method: POST,
            requestContentType: ContentType.JSON,
            body: [[
               ref: "refs/heads/$branchName".toString(),
               sha: "$shaRef".toString()
            ], null]).data
    }

    static exchangeOAuthToken(String code) {
        final client = new HttpURLClient(url: 'https://github.com/login/oauth/access_token')
        client.request(query: [
                client_id: Credential.githubClientId,
                client_secret: Credential.githubClientSecret,
                code: code
        ]).data
    }

}

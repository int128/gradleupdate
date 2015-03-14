package infrastructure

import groovyx.net.http.ContentType
import groovyx.net.http.HttpResponseException
import groovyx.net.http.HttpURLClient

import static groovyx.net.http.Method.DELETE
import static groovyx.net.http.Method.POST

class GitHubRepository {

    private final String repo

    private final HttpURLClient client

    def GitHubRepository(String repo, String token) {
        this.repo = repo
        this.client = new HttpURLClient(url: 'https://api.github.com', headers: [
                'Authorization': "token $token",
                'User-Agent': 'gradleupdate'
        ])
    }

    def createBranch(String branchName, String from) {
        def sha = getReference(from).object.sha
        assert sha
        createReference(branchName, sha)
    }

    boolean removeBranch(String branchName) {
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

    def getReference(String branchName) {
        client.request(path: "/repos/$repo/git/refs/heads/$branchName").data
    }

    def createReference(String branchName, String shaRef) {
        client.request(path: "/repos/$repo/git/refs", method: POST,
            requestContentType: ContentType.JSON,
            body: [[
               ref: "refs/heads/$branchName".toString(),
               sha: "$shaRef".toString()
            ], null]).data
    }
}

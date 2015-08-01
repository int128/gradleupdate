package infrastructure

import groovyx.net.http.ContentType
import groovyx.net.http.HttpResponseException
import groovyx.net.http.HttpURLClient
import model.Credential

import static groovyx.net.http.Method.DELETE
import static groovyx.net.http.Method.POST

class GitHub {

    private final HttpURLClient client

    private final Credential credential

    def GitHub(Map headers = [:]) {
        credential = Credential.getOrCreate('github')
        client = new HttpURLClient(url: 'https://api.github.com', headers: [
                'Authorization': "token $credential.token",
                'User-Agent': 'gradleupdate'
        ] + headers)
    }

    static GitHub authorizationOrDefault(String authorizationHeader) {
        if (authorizationHeader) {
            assert authorizationHeader.startsWith('token ')
            new GitHub(Authorization: authorizationHeader)
        } else {
            new GitHub()
        }
    }

    def getRepository(String repo) {
        client.request(path: "/repos/$repo").data
    }

    void deleteRepository(String repo) {
        try {
            client.request(path: "/repos/$repo", method: DELETE).data
        } catch (NullPointerException ignore) {
            // 204 No Content causes NPE due to the bug of HttpURLClient
        }
    }

    def fork(String repo) {
        client.request(path: "/repos/$repo/forks", method: POST).data
    }

    def getRepositories(String userName) {
        client.request(path: "/users/$userName/repos").data
    }

    def getStargazers(String repo) {
        client.request(path: "/repos/$repo/stargazers").data
    }

    def getContent(String repo, String path) {
        client.request(path: "/repos/$repo/contents/$path").data
    }

    def createBranch(String repo, String branchName, String from) {
        def sha = getReference(repo, from).object.sha
        assert sha
        createReference(repo, branchName, sha)
    }

    def createBranch(String repo, String branchName, String from, String message, List<Map> contents) {
        def ref = getReference(repo, from).object.sha
        assert ref instanceof String
        def tree = getCommit(repo, ref).tree.sha
        assert tree instanceof String
        def newTree = createTree(repo, tree, contents).sha
        assert newTree instanceof String
        def newCommit = createCommit(repo, [ref], newTree, message).sha
        assert newCommit instanceof String
        createReference(repo, branchName, newCommit)
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

    def createPullRequest(String repo, String into, String from, String title, String body) {
        requestJson(path: "/repos/$repo/pulls", method: POST, body: [
                head: from, base: into, title: title, body: body
        ]).data
    }

    def getReference(String repo, String branchName) {
        client.request(path: "/repos/$repo/git/refs/heads/$branchName").data
    }

    def createReference(String repo, String branchName, String shaRef) {
        requestJson(path: "/repos/$repo/git/refs", method: POST, body: [
                ref: "refs/heads/$branchName".toString(), sha: "$shaRef".toString()
        ]).data
    }

    def getCommit(String repo, String sha) {
        client.request(path: "/repos/$repo/git/commits/$sha").data
    }

    def createCommit(String repo, List<String> parents, String tree, String message) {
        requestJson(path: "/repos/$repo/git/commits", method: POST, body: [
                parents: parents, tree: tree, message: message
        ]).data
    }

    def createTree(String repo, String baseSha, List<Map> contents) {
        requestJson(path: "/repos/$repo/git/trees", method: POST, body: [base_tree: baseSha, tree: contents]).data
    }

    def createBlob(String repo, String content, String encoding = 'base64') {
        requestJson(path: "/repos/$repo/git/blobs", method: POST, body: [content: content, encoding: encoding]).data
    }

    def exchangeOAuthToken(String code) {
        final client = new HttpURLClient(url: 'https://github.com/login/oauth/access_token')
        client.request(query: [
                client_id: credential.clientId,
                client_secret: credential.clientSecret,
                code: code
        ]).data
    }

    private requestJson(Map request) {
        client.request(request + [requestContentType: ContentType.JSON, body: [request.body, null]])
    }

}

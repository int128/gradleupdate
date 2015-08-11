package infrastructure

import groovy.util.logging.Log
import groovyx.net.http.ContentType
import groovyx.net.http.HttpResponseException
import groovyx.net.http.HttpURLClient
import model.Credential

import static groovyx.net.http.Method.DELETE
import static groovyx.net.http.Method.POST

@Log
class GitHub {

    private final HttpURLClient client

    def GitHub(Map headers = [:]) {
        def token = Credential.getOrCreate('github-token')
        client = new HttpURLClient(url: 'https://api.github.com', headers: [
                'Authorization': "token $token.secret",
                'User-Agent': 'gradleupdate'
        ] + headers)
    }

    boolean removeRepository(String repo) {
        handleHttpResponseException(404: false) {
            handle204NoContentWorkaround(true) {
                // 204 No Content causes NPE due to the bug of HttpURLClient
                client.request(path: "/repos/$repo", method: DELETE).data
            }
        }
    }

    def fork(String repo) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/forks", method: POST).data
        }
    }

    def fetchRepositories(String userName) {
        handleHttpResponseException(404: null) {
            client.request(path: "/users/$userName/repos").data
        }
    }

    def fetchStargazers(String repo) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/stargazers").data
        }
    }

    def fetchContent(String repo, String path) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/contents/$path").data
        }
    }

    def createBranch(String repo, String branchName, String from) {
        handleHttpResponseException(404: null) {
            def sha = fetchReference(repo, from).object.sha
            assert sha instanceof String
            createReference(repo, branchName, sha)
        }
    }

    def createBranch(String repo, String branchName, String from, String message, List<Map> contents) {
        def ref = fetchReference(repo, from).object.sha
        assert ref instanceof String
        def tree = fetchCommit(repo, ref).tree.sha
        assert tree instanceof String
        def newTree = createTree(repo, tree, contents).sha
        assert newTree instanceof String
        def newCommit = createCommit(repo, [ref], newTree, message).sha
        assert newCommit instanceof String
        createReference(repo, branchName, newCommit)
    }

    boolean removeBranch(String repo, String branchName) {
        // API returns 422 if branch does not exist
        handleHttpResponseException(422: false) {
            // 204 No Content causes NPE due to the bug of HttpURLClient
            handle204NoContentWorkaround(true) {
                client.request(path: "/repos/$repo/git/refs/heads/$branchName", method: DELETE).success
            }
        }
    }

    def fetchPullRequests(Map filter, String repo) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/pulls", query: filter).data
        }
    }

    def createPullRequest(String repo, String into, String from, String title, String body) {
        handleHttpResponseException(404: null) {
            requestJson(path: "/repos/$repo/pulls", method: POST, body: [
                    head: from, base: into, title: title, body: body
            ]).data
        }
    }

    def fetchReference(String repo, String branchName) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/git/refs/heads/$branchName").data
        }
    }

    def createReference(String repo, String branchName, String shaRef) {
        handleHttpResponseException(404: null) {
            requestJson(path: "/repos/$repo/git/refs", method: POST, body: [
                    ref: "refs/heads/$branchName".toString(), sha: "$shaRef".toString()
            ]).data
        }
    }

    def fetchCommit(String repo, String sha) {
        handleHttpResponseException(404: null) {
            client.request(path: "/repos/$repo/git/commits/$sha").data
        }
    }

    def createCommit(String repo, List<String> parents, String tree, String message) {
        handleHttpResponseException(404: null) {
            requestJson(path: "/repos/$repo/git/commits", method: POST, body: [
                    parents: parents, tree: tree, message: message
            ]).data
        }
    }

    def createTree(String repo, String baseSha, List<Map> contents) {
        handleHttpResponseException(404: null) {
            requestJson(path: "/repos/$repo/git/trees", method: POST, body: [base_tree: baseSha, tree: contents]).data
        }
    }

    def createBlob(String repo, String content, String encoding = 'base64') {
        handleHttpResponseException(404: null) {
            requestJson(path: "/repos/$repo/git/blobs", method: POST, body: [content: content, encoding: encoding]).data
        }
    }

    private static handleHttpResponseException(Map statusCodeMap, Closure closure) {
        try {
            closure()
        } catch (HttpResponseException e) {
            def statusCode = e.response.status
            if (statusCodeMap.containsKey(statusCode)) {
                def value = statusCodeMap[statusCode]
                log.warning("Got error response $statusCode and returning $value")
                value
            } else {
                throw e
            }
        }
    }

    private static handle204NoContentWorkaround(Object value, Closure closure) {
        try {
            closure()
        } catch (NullPointerException e) {
            log.info("204 No Content caused NPE but ignored: $e.localizedMessage")
            value
        }
    }

    private requestJson(Map request) {
        client.request(request + [requestContentType: ContentType.JSON, body: [request.body, null]])
    }

}

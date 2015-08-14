package infrastructure

import groovy.util.logging.Log
import model.Credential
import wslite.rest.ContentType
import wslite.rest.RESTClient

@Log
class GitHub implements ErrorStatusHandler {

    private final RESTClient client

    def GitHub() {
        this(Credential.getOrCreate('github-token'))
    }

    def GitHub(Credential credential) {
        def headers = ['User-Agent': 'gradleupdate']
        if (credential) {
            headers += [Authorization: "token $credential.secret"]
        }
        client = new RESTClient('https://api.github.com')
        client.httpClient.defaultHeaders += headers
    }

    boolean removeRepository(String repo) {
        handleHttpResponseException(404: false) {
            client.delete(path: "/repos/$repo")
            true
        }
    }

    def fork(String repo) {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/forks").json
        }
    }

    def fetchRepositories(Map filter = [:], String userName) {
        handleHttpResponseException(404: null) {
            // TODO: fetch from more pages
            client.get(path: "/users/$userName/repos", query: [per_page: 100] + filter).json
        }
    }

    def fetchStargazers(String repo) {
        handleHttpResponseException(404: null) {
            // TODO: fetch from more pages
            client.get(path: "/repos/$repo/stargazers", query: [per_page: 100]).json
        }
    }

    def fetchContent(String repo, String path) {
        handleHttpResponseException(404: null) {
            client.get(path: "/repos/$repo/contents/$path").json
        }
    }

    def createBranch(String repo, String branchName, String from) {
        // TODO: should move to domain class
        handleHttpResponseException(404: null) {
            def sha = fetchReference(repo, from).object.sha
            assert sha instanceof String
            createReference(repo, branchName, sha)
        }
    }

    boolean removeBranch(String repo, String branchName) {
        // API returns 422 if branch does not exist
        handleHttpResponseException(422: false) {
            client.delete(path: "/repos/$repo/git/refs/heads/$branchName")
            true
        }
    }

    def fetchPullRequests(Map filter, String repo) {
        handleHttpResponseException(404: null) {
            client.get(path: "/repos/$repo/pulls", query: filter).json
        }
    }

    def createPullRequest(String repo, String base, String head, String title, String body) {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/pulls") {
                type ContentType.JSON
                json head: head, base: base, title: title, body: body
            }.json
        }
    }

    def fetchReference(String repo, String branchName) {
        handleHttpResponseException(404: null) {
            client.get(path: "/repos/$repo/git/refs/heads/$branchName").json
        }
    }

    def createReference(String repo, String branchName, String shaRef) {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/git/refs") {
                type ContentType.JSON
                json ref: "refs/heads/$branchName", sha: shaRef
            }.json
        }
    }

    def fetchCommit(String repo, String sha) {
        handleHttpResponseException(404: null) {
            client.get(path: "/repos/$repo/git/commits/$sha").json
        }
    }

    def createCommit(String repo, List<String> parents, String tree, String message) {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/git/commits") {
                type ContentType.JSON
                json parents: parents, tree: tree, message: message
            }.json
        }
    }

    def createTree(String repo, String baseSha, List<Map> contents) {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/git/trees") {
                type ContentType.JSON
                json base_tree: baseSha, tree: contents
            }.json
        }
    }

    def createBlob(String repo, String content, String encoding = 'base64') {
        handleHttpResponseException(404: null) {
            client.post(path: "/repos/$repo/git/blobs") {
                type ContentType.JSON
                json content: content, encoding: encoding
            }.json
        }
    }

}

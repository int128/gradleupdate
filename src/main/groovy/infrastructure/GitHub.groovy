package infrastructure

import groovy.util.logging.Log
import model.Credential
import wslite.rest.ContentType
import wslite.rest.RESTClient

import static model.Credential.CredentialKey.GitHubToken

@Log
class GitHub implements ErrorStatusHandler {

    private final RESTClient client

    def GitHub() {
        this(Credential.get(GitHubToken))
    }

    def GitHub(Credential credential) {
        def headers = ['User-Agent': 'gradleupdate']
        if (credential) {
            headers += [Authorization: "token $credential.secret"]
        }
        client = new RESTClient('https://api.github.com')
        client.httpClient.defaultHeaders += headers
    }

    def fetch(String repo) {
        handleHttpResponseException(404: false) {
            client.get(path: "/repos/$repo").json
        }
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

    def fetchRepositoriesOfFirstPage(Map filter = [:], String userName) {
        handleHttpResponseException(404: null) {
            new Page(client.get(path: "/users/$userName/repos", query: filter))
        } as Page
    }

    def fetchStargazersOfFirstPage(String repo) {
        handleHttpResponseException(404: null) {
            new Page(client.get(path: "/repos/$repo/stargazers"))
        } as Page
    }

    def fetchNextPage(String relation) {
        assert relation
        assert relation.startsWith(client.url)
        def path = relation.substring(client.url.length())
        assert path.startsWith('/')
        handleHttpResponseException(404: null) {
            new Page(client.get(path: path))
        } as Page
    }

    def fetchContent(String repo, String path) {
        handleHttpResponseException(404: null) {
            client.get(path: "/repos/$repo/contents/$path").json
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

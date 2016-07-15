package domain

import groovy.util.logging.Log
import wslite.rest.ContentType

@Log
class GHCommit {

    final GHRepository repository
    final def rawJson

    private def GHCommit(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    static GHCommit get(GHBranch branch) {
        assert branch
        log.info("Fetching commit of branch $branch")
        def repository = branch.repository
        def response = repository.session.client.get(path: "/repos/$repository/git/commits/$branch.sha.value")
        assert response.statusCode == 200
        new GHCommit(repository, response.json)
    }

    static GHCommit create(GHBranch branch, String message, List<GHTreeContent> contents) {
        assert branch
        assert message
        def tree = get(branch).tree.upload(contents)
        log.info("Creating commit on branch $branch with message $message")
        def repository = branch.repository
        def response = repository.session.client.post(path: "/repos/$repository/git/commits") {
            type ContentType.JSON
            json parents: [branch.sha.value], tree: tree.sha.value, message: message
        }
        assert response.statusCode == 201
        new GHCommit(repository, response.json)
    }

    @Lazy
    GHCommitSha sha = { new GHCommitSha(rawJson.sha as String) }()

    @Lazy
    GHTree tree = {
        assert rawJson.tree.sha instanceof String
        new GHTree(repository, new GHTreeSha(rawJson.tree.sha as String))
    }()

    @Override
    String toString() {
        "$repository/$sha"
    }

}

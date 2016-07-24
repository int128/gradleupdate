package domain

import groovy.util.logging.Log

@Log
class GHBranch {

    final GHRepository repository
    final def rawJson

    def GHBranch(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    @Lazy
    String name = {
        assert rawJson.ref.startsWith('refs/heads/')
        rawJson.ref.substring('refs/heads/'.length())
    }()

    @Lazy
    GHCommitSha sha = {
        assert rawJson.object.type == 'commit'
        new GHCommitSha(rawJson.object.sha as String)
    }()

    GHBranch syncTo(GHCommitSha newSha) {
        assert newSha
        sha == newSha ? this : repository.refs.update("heads/$name", sha: newSha.value, force: true)
    }

    GHBranch clone(String intoName) {
        assert intoName
        repository.createBranch(intoName, sha)
    }

    GHCommit commit(String message, List<GHTreeContent> contents) {
        assert message
        assert contents
        def baseTree = repository.commits.get(sha.value).tree
        def tree = repository.trees.create(
                base_tree: baseTree.sha.value,
                tree: contents.collect { content ->
                    def blob = repository.blobs.create(encoding: 'base64', content: content.base64encoded)
                    new GHTreeFile(content.path, content.mode, blob.sha).asMap()
                })
        repository.commits.create(parents: [sha.value], tree: tree.sha.value, message: message)
    }

    @Override
    String toString() {
        "GHBranch($name:$sha)@$repository"
    }

    GHContent getContent(String path) {
        repository.contents.get(path, ref: name)
    }

    GHContent findContent(String path) {
        repository.contents.find(path, ref: name)
    }

}

package domain

class GHCommit {

    final GHRepository repository
    final def rawJson

    def GHCommit(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    @Lazy
    GHCommitSha sha = { new GHCommitSha(rawJson.sha as String) }()

    @Lazy
    GHTree tree = {
        assert rawJson.tree.sha instanceof String
        new GHTree(repository, rawJson.tree)
    }()

    @Override
    String toString() {
        "$repository/commits/$sha"
    }

}

package domain

class GHTree {

    final GHRepository repository
    final def rawJson

    def GHTree(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    @Lazy
    GHTreeSha sha = {
        assert rawJson.sha
        new GHTreeSha(rawJson.sha as String)
    }()

    @Override
    String toString() {
        "$repository/trees/$sha"
    }

}

package domain

class GHBlob {

    final GHRepository repository
    final def rawJson

    def GHBlob(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    @Lazy
    GHBlobSha sha = {
        assert rawJson.sha
        new GHBlobSha(rawJson.sha as String)
    }()

    @Override
    String toString() {
        "$repository/blobs/$sha"
    }

}

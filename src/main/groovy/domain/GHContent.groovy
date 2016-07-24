package domain

class GHContent {

    final GHRepository repository
    final def rawJson

    def GHContent(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    @Lazy
    String path = { rawJson.path as String }()

    @Lazy
    String base64encoded = {
        assert rawJson.encoding == 'base64'
        rawJson.content as String
    }()

    @Lazy
    byte[] data = { base64encoded.decodeBase64() }()

    @Lazy
    String contentAsString = { new String(data) }()

    @Override
    String toString() {
        "GHContent($path)@$repository"
    }

}

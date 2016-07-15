package domain

import groovy.util.logging.Log

@Log
class GHContent {

    final def rawJson

    def GHContent(def rawJson) {
        this.rawJson = rawJson
    }

    static GHContent get(GHBranch branch, String path) {
        assert branch
        assert path
        log.info("Fetching content from branch $branch: $path")
        def response = branch.repository.session.client.get(
                path: "/repos/$branch.repository/contents/$path",
                query: [ref: branch.name]
        )
        assert response.statusCode == 200
        new GHContent(response.json)
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

}

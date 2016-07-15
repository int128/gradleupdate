package domain

import groovy.util.logging.Log
import wslite.rest.ContentType

@Log
class GHBlob {

    final GHRepository repository
    final GHBlobSha sha

    private def GHBlob(GHRepository repository, GHBlobSha sha) {
        this.repository = repository
        this.sha = sha
    }

    static GHBlob create(GHRepository repository, String encoded) {
        assert repository
        log.info("Uploading encoded ${encoded.length()} bytes into $repository")
        def response = repository.session.client.post(path: "/repos/$repository/git/blobs") {
            type ContentType.JSON
            json encoding: 'base64', content: encoded
        }
        assert response.statusCode == 201
        assert response.json.sha instanceof String
        new GHBlob(repository, new GHBlobSha(response.json.sha as String))
    }

    @Override
    String toString() {
        "$repository/$sha"
    }

}

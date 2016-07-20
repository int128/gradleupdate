package domain

import groovy.util.logging.Log
import wslite.rest.ContentType

@Log
class GHTree {

    final GHRepository repository
    final GHTreeSha sha

    def GHTree(GHRepository repository, GHTreeSha sha) {
        this.repository = repository
        this.sha = sha
    }

    GHTree upload(List<GHTreeContent> contents) {
        assert contents
        create(contents.collect { content ->
            def blob = GHBlob.create(repository, content.base64encoded)
            new GHTreeFile(content.path, content.mode, blob.sha)
        })
    }

    GHTree create(List<GHTreeFile> files) {
        assert files
        log.info("Creating tree with ${files.size()} files on $this")
        def created = repository.session.client.post(path: "/repos/$repository/git/trees") {
            type ContentType.JSON
            json base_tree: sha.value, tree: files*.asMap()
        }
        assert created.statusCode == 201
        assert created.json.sha instanceof String
        new GHTree(repository, new GHTreeSha(created.json.sha as String))
    }

    @Override
    String toString() {
        "$repository/$sha"
    }

}

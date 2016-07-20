package domain

import groovy.util.logging.Log
import wslite.rest.ContentType

@Log
class GHBranch {

    final GHRepository repository
    final def rawJson

    private def GHBranch(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    static GHBranch get(GHRepository repository, String name) {
        assert repository
        assert name
        log.info("Fetching branch $name of repository $repository")
        def response = repository.session.client.get(path: "/repos/$repository/git/refs/heads/$name")
        assert response.statusCode == 200
        new GHBranch(repository, response.json)
    }

    static GHBranch create(GHRepository repository, String name, GHCommitSha sha) {
        assert repository
        assert name
        log.info("Creating branch $name on repository $repository")
        def response = repository.session.client.post(path: "/repos/$repository/git/refs") {
            type ContentType.JSON
            json ref: "refs/heads/$name", sha: sha.value
        }
        assert response.statusCode == 201
        new GHBranch(repository, response.json)
    }

    static GHBranch update(GHRepository repository, String name, GHCommitSha sha, boolean force) {
        log.info("Updating branch $name of repository $repository to $sha")
        def response = repository.session.client.patch(path: "/repos/$repository/git/refs/heads/$name") {
            type ContentType.JSON
            json sha: sha.value, force: force
        }
        assert response.statusCode == 200
        new GHBranch(repository, response.json)
    }

    static boolean remove(GHRepository repository, String name) {
        assert repository
        assert name
        log.info("Removing branch $name of repository $repository")
        def response = repository.session.client.delete(path: "/repos/$repository/git/refs/heads/$name")
        assert response.statusCode in [204, 422]
        // API returns 422 if branch does not exist
        response.statusCode == 204
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
        if (sha == newSha) {
            this
        } else {
            update(repository, name, newSha, true)
        }
    }

    GHBranch clone(String intoName) {
        log.info("Cloning branch $name of repository $repository as $intoName")
        create(repository, intoName, sha)
    }

    GHCommit commit(String message, List<GHTreeContent> contents) {
        GHCommit.create(this, message, contents)
    }

    @Override
    String toString() {
        "$repository:$name"
    }

}

package domain

import groovy.util.logging.Log

import static infrastructure.HTTPClientExceptionUtil.nullIfResourceIsNotFound

@Log
class GHRepository {

    final GHSession session
    final def rawJson

    private def GHRepository(GHSession session, def rawJson) {
        this.session = session
        this.rawJson = rawJson
    }

    static GHRepository get(GHSession session, String fullName) {
        assert session
        assert fullName
        log.info("Fetching repository $fullName")
        def response = session.client.get(path: "/repos/$fullName")
        assert response.statusCode == 200
        new GHRepository(session, response.json)
    }

    @Lazy
    String fullName = { rawJson.full_name }()

    @Lazy
    String ownerName = { rawJson.owner.login }()

    GHRepository fork() {
        log.info("Creating fork from repository $this")
        def response = session.client.post(path: "/repos/$fullName/forks")
        assert response.statusCode == 202
        new GHRepository(session, response.json)
    }

    boolean checkPermissionPush() {
        rawJson.permissions?.push == true
    }

    boolean checkPermissionPull() {
        rawJson.permissions?.pull == true
    }

    GHBranch getBranch(String name) {
        GHBranch.get(this, name)
    }

    @Lazy
    GHBranch defaultBranch = {
        assert rawJson.default_branch
        GHBranch.get(this, rawJson.default_branch)
    }()

    GHBranch createBranch(String name, GHCommitSha reference) {
        GHBranch.create(this, name, reference)
    }

    boolean removeBranch(String name) {
        GHBranch.remove(this, name)
    }

    GHBranch createOrResetBranch(String name, GHCommitSha sha) {
        nullIfResourceIsNotFound {
            GHBranch.update(this, name, sha, true)
        } ?: createBranch(name, sha)
    }

    @Override
    String toString() {
        fullName
    }

}

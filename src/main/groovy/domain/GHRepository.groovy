package domain

import infrastructure.RestAPI

class GHRepository {

    final GHSession session
    final def rawJson

    def GHRepository(GHSession session, def rawJson) {
        this.session = session
        this.rawJson = rawJson
    }

    @Lazy
    String fullName = { rawJson.full_name }()

    @Lazy
    String ownerName = { rawJson.owner.login }()

    @Lazy
    def forks = { RestAPI.of(GHRepository, "/repos/$fullName/forks", session, session.client) }()

    @Lazy
    def refs = { RestAPI.of(GHBranch, "/repos/$fullName/git/refs", this, session.client) }()

    @Lazy
    def commits = { RestAPI.of(GHCommit, "/repos/$fullName/git/commits", this, session.client) }()

    @Lazy
    def trees = { RestAPI.of(GHTree, "/repos/$fullName/git/trees", this, session.client) }()

    @Lazy
    def blobs = { RestAPI.of(GHBlob, "/repos/$fullName/git/blobs", this, session.client) }()

    @Lazy
    def contents = { RestAPI.of(GHContent, "/repos/$fullName/contents", this, session.client) }()

    @Lazy
    def pullRequests = { RestAPI.of(GHPullRequest, "/repos/$fullName/pulls", this, session.client) }()

    GHRepository fork() {
        forks.invoke()
    }

    boolean checkPermissionPush() {
        rawJson.permissions?.push == true
    }

    GHBranch getBranch(String name) {
        refs.get("heads/$name")
    }

    @Lazy
    GHBranch defaultBranch = {
        assert rawJson.default_branch
        getBranch(rawJson.default_branch)
    }()

    GHBranch createBranch(String name, GHCommitSha sha) {
        refs.create(ref: "refs/heads/$name", sha: sha.value)
    }

    boolean removeBranch(String name) {
        refs.delete("heads/$name")
    }

    GHBranch createOrResetBranch(String name, GHCommitSha sha) {
        refs.find("heads/$name")?.syncTo(sha) ?: createBranch(name, sha)
    }

    @Override
    String toString() {
        "GHRepository($fullName)"
    }

}

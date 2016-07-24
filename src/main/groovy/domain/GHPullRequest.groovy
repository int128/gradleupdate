package domain

class GHPullRequest {

    final GHRepository repository
    final def rawJson

    def GHPullRequest(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    static GHPullRequest findLatest(GHBranch forkBranch, GHBranch originBranch) {
        firstOrNull(originBranch.repository.pullRequests.findAll(
                head: "${forkBranch.repository.ownerName}:${forkBranch.name}",
                base: originBranch.name,
                state: 'all',
                direction: 'desc',
        ))
    }

    static GHPullRequest create(GHBranch forkBranch, GHBranch originBranch, String title, String body) {
        originBranch.repository.pullRequests.create(
                head: "${forkBranch.repository.ownerName}:${forkBranch.name}",
                base: originBranch.name,
                title: title,
                body: body
        )
    }

    static GHPullRequest update(GHRepository repository, int number, String title, String body) {
        repository.pullRequests.update(number as String, title: title, body: body)
    }

    static GHPullRequest createOrUpdate(GHBranch forkBranch, GHBranch originBranch, String title, String body) {
        findLatest(forkBranch, originBranch)?.update(title, body) ?: create(forkBranch, originBranch, title, body)
    }

    @Lazy
    int number = { rawJson.number }()

    @Lazy
    String url = { rawJson.html_url }()

    @Lazy
    Date createdAt = { Date.parse("yyyy-MM-dd'T'HH:mm:ss", rawJson.created_at) }()

    GHPullRequest update(String title, String body) {
        update(repository, number, title, body)
    }

    @Override
    String toString() {
        "GHPullRequest($number)@$repository"
    }

    private static <E> E firstOrNull(List<E> list) {
        list.empty ? null : list.first()
    }

}

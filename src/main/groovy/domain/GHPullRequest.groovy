package domain

import groovy.util.logging.Log
import wslite.rest.ContentType

@Log
class GHPullRequest {

    final GHRepository repository
    final def rawJson

    def GHPullRequest(GHRepository repository, def rawJson) {
        this.repository = repository
        this.rawJson = rawJson
    }

    static List<GHPullRequest> find(GHBranch forkBranch, GHBranch originBranch) {
        log.info("Finding pull request: $forkBranch -> $originBranch")
        def response = originBranch.repository.session.client.get(path: "/repos/$originBranch.repository/pulls", query: [
                head: "${forkBranch.repository.ownerName}:${forkBranch.name}",
                base: originBranch.name,
        ])
        assert response.statusCode == 200
        response.json.collect { single ->
            new GHPullRequest(originBranch.repository, single)
        }
    }

    static GHPullRequest create(GHBranch forkBranch, GHBranch originBranch, String title, String body) {
        log.info("Creating pull request: $forkBranch -> $originBranch")
        def response = originBranch.repository.session.client.post(path: "/repos/$originBranch.repository/pulls") {
            type ContentType.JSON
            json head: "${forkBranch.repository.ownerName}:${forkBranch.name}", base: originBranch.name, title: title, body: body
        }
        assert response.statusCode == 201
        new GHPullRequest(forkBranch.repository, response.json)
    }

    static GHPullRequest update(GHRepository repository, int number, String title, String body) {
        log.info("Updating pull request: $repository/PullRequest($number)")
        def response = repository.session.client.patch(path: "/repos/$repository/pulls/$number") {
            type ContentType.JSON
            json title: title, body: body
        }
        assert response.statusCode == 200
        new GHPullRequest(repository, response.json)
    }

    static GHPullRequest createOrUpdate(GHBranch forkBranch, GHBranch originBranch, String title, String body) {
        def existing = find(forkBranch, originBranch)
        if (existing.empty) {
            create(forkBranch, originBranch, title, body)
        } else {
            existing.last().update(title, body)
        }
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
        "$repository/PullRequest($number)"
    }

}

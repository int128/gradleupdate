package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.PullRequest
import org.eclipse.egit.github.core.Repository
import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.PullRequestService

class EnhancedPullRequestService(client: GitHubClient) : PullRequestService(client) {
    class Query(
        val state: String = "",
        val head: String = "",
        val base: String = "",
        val sort: String = "",
        val direction: String = "",
        val start: Int = 1,
        val size: Int
    )

    /**
     * @see https://developer.github.com/v3/pulls/#list-pull-requests
     */
    fun query(repository: Repository, query: Query): List<PullRequest> {
        val request = createPullsRequest(repository, null, query.start, query.size)
        request.params = mapOf(
            "state" to query.state,
            "head" to query.head,
            "base" to query.base,
            "sort" to query.sort,
            "direction" to query.direction
        ).filter { it.value.isNotEmpty() }
        return getAll(request)
    }
}

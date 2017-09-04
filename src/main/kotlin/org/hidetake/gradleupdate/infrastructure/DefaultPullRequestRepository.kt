package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.*
import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.DataService
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.GradleWrapperPullRequest
import org.hidetake.gradleupdate.domain.PullRequestRepository
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedPullRequestService
import java.util.*

@org.springframework.stereotype.Repository
class DefaultPullRequestRepository(client: GitHubClient) : PullRequestRepository {
    private val repositoryService = RepositoryService(client)
    private val pullRequestService = EnhancedPullRequestService(client)
    private val dataService = DataService(client)

    override fun createOrUpdate(gradleWrapperPullRequest: GradleWrapperPullRequest): PullRequest {
        val baseRepository = repositoryService.getRepository({gradleWrapperPullRequest.repositoryName})
        val baseCommit = dataService.getCommit(baseRepository,
            dataService.getReference(baseRepository,
                "refs/heads/${baseRepository.defaultBranch}").`object`.sha)
        val fork = repositoryService.forkRepository(baseRepository)
        createOrUpdateBranchWithCommit(fork, baseCommit, gradleWrapperPullRequest)
        return createOrUpdatePullRequest(baseRepository, fork, gradleWrapperPullRequest)
    }

    private fun createOrUpdatePullRequest(
        baseRepository: Repository,
        headRepository: Repository,
        gradleWrapperPullRequest: GradleWrapperPullRequest
    ): PullRequest {
        val query = EnhancedPullRequestService.Query(
            base = baseRepository.defaultBranch,
            head = "${headRepository.owner.login}:${gradleWrapperPullRequest.branchName}",
            state = "open",
            start = 1,
            size = 1
        )
        val latest = pullRequestService.query(baseRepository, query).firstOrNull()
        return if (latest == null) {
            pullRequestService.createPullRequest(baseRepository, PullRequest().apply {
                title = gradleWrapperPullRequest.title
                body = gradleWrapperPullRequest.description
                base = PullRequestMarker().apply { label = query.base }
                head = PullRequestMarker().apply { label = query.head }
            })
        } else {
            pullRequestService.editPullRequest(baseRepository, latest.apply {
                title = gradleWrapperPullRequest.title
                body = gradleWrapperPullRequest.description
            })
        }
    }

    private fun createOrUpdateBranchWithCommit(
        repository: Repository,
        parent: Commit,
        gradleWrapperPullRequest: GradleWrapperPullRequest
    ): Reference {
        val refName = "refs/heads/${gradleWrapperPullRequest.branchName}"
        val existentRef = nullIfNotFound {
            dataService.getReference(repository, refName)
        }
        return if (existentRef == null) {
            dataService.createReference(repository, Reference().apply {
                ref = refName
                `object` = TypedResource().apply {
                    sha = createCommit(repository, parent, gradleWrapperPullRequest).sha
                }
            })
        } else {
            val existentRefParent = dataService.getCommit(repository, existentRef.`object`.sha).parents.firstOrNull()
            val needToUpdateBranch = existentRefParent?.sha == parent.sha
            if (!needToUpdateBranch) {
                dataService.editReference(repository, Reference().apply {
                    ref = refName
                    `object` = TypedResource().apply {
                        sha = createCommit(repository, parent, gradleWrapperPullRequest).sha
                    }
                }, true)
            } else {
                existentRef
            }
        }
    }

    private fun createCommit(repository: Repository, parent: Commit, gradleWrapperPullRequest: GradleWrapperPullRequest) =
        dataService.createCommit(repository, Commit().apply {
            author = CommitUser().apply {
                name = gradleWrapperPullRequest.authorName
                email = gradleWrapperPullRequest.authorEmail
                date = Date()
            }
            committer = author
            message = gradleWrapperPullRequest.title
            parents = listOf(parent)
            tree = dataService.createTree(repository,
                gradleWrapperPullRequest.files.map { file ->
                    TreeEntry().apply {
                        path = file.path
                        mode = when (file.executable) {
                            true  -> TreeEntry.MODE_BLOB_EXECUTABLE
                            false -> TreeEntry.MODE_BLOB
                        }
                        sha = dataService.createBlob(repository, Blob().apply {
                            content = file.base64Content
                            encoding = Blob.ENCODING_BASE64
                        })
                    }
                },
                parent.tree.sha)
        })
}

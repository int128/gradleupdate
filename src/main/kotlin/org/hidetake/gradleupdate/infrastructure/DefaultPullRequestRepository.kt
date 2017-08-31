package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.*
import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.DataService
import org.eclipse.egit.github.core.service.PullRequestService
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.GradleWrapperPullRequest
import org.hidetake.gradleupdate.domain.PullRequestRepository

@org.springframework.stereotype.Repository
class DefaultPullRequestRepository(client: GitHubClient) : PullRequestRepository {
    private val repositoryService = RepositoryService(client)
    private val pullRequestService = PullRequestService(client)
    private val dataService = DataService(client)

    override fun create(repositoryName: String, gradleWrapperPullRequest: GradleWrapperPullRequest) {
        val headBranchName = gradleWrapperPullRequest.branchName

        val baseRepository = repositoryService.getRepository({repositoryName})
        val baseRef = dataService.getReference(baseRepository, "refs/heads/${baseRepository.defaultBranch}")

        val fork = repositoryService.forkRepository(baseRepository)

        val headRef = nullIfNotFound {dataService.getReference(fork, "refs/heads/$headBranchName")}
        when (headRef) {
            null ->
                dataService.createReference(fork, Reference().apply {
                    ref = "refs/heads/$headBranchName"
                    `object` = TypedResource().apply {
                        sha = createCommit(fork, baseRef, gradleWrapperPullRequest).sha
                    }
                })

            else ->
                when (dataService.getCommit(fork, headRef.`object`.sha).parents.firstOrNull()?.sha) {
                    baseRef.`object`.sha -> headRef

                    else ->
                        dataService.editReference(fork, Reference().apply {
                            ref = "refs/heads/$headBranchName"
                            `object` = TypedResource().apply {
                                sha = createCommit(fork, baseRef, gradleWrapperPullRequest).sha
                            }
                        }, true)
                }
        }

        pullRequestService.createPullRequest({repositoryName}, PullRequest().apply {
            title = gradleWrapperPullRequest.title
            body = gradleWrapperPullRequest.description
            base = PullRequestMarker().apply { label = headBranchName }
            head = PullRequestMarker().apply { label = "${fork.owner.login}:$headBranchName" }
        })
    }

    private fun createCommit(repository: Repository, parent: Reference, gradleWrapperPullRequest: GradleWrapperPullRequest) =
        dataService.createCommit(repository, Commit().apply {
            message = gradleWrapperPullRequest.title
            parents = listOf(Commit().apply { sha = parent.`object`.sha })
            tree = dataService.createTree(repository, gradleWrapperPullRequest.files.map { file ->
                TreeEntry().apply {
                    path = file.path
                    mode = if (file.executable) TreeEntry.MODE_BLOB_EXECUTABLE else TreeEntry.MODE_BLOB
                    sha = dataService.createBlob(repository, Blob().apply {
                        content = file.base64Content
                        encoding = Blob.ENCODING_BASE64
                    })
                }
            })
        })}

package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.*
import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.DataService
import org.eclipse.egit.github.core.service.RepositoryService
import org.eclipse.egit.github.core.service.UserService
import org.hidetake.gradleupdate.domain.*
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedPullRequestService

@org.springframework.stereotype.Repository
class DefaultPullRequestRepository(client: GitHubClient) : PullRequestRepository {
    private val userService = UserService(client)
    private val repositoryService = RepositoryService(client)
    private val pullRequestService = EnhancedPullRequestService(client)
    private val dataService = DataService(client)

    override fun createOrUpdate(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion,
        files: List<GradleWrapperFile>
    ) {
        val branch = Branch("gradle-${gradleWrapperVersion.version}-${repositoryPath.owner}")
        val title = "Gradle ${gradleWrapperVersion.version}"
        val description = "Gradle ${gradleWrapperVersion.version} is available."

        val baseRepository = repositoryService.getRepository({repositoryPath.fullName})
        val baseRef = dataService.getReference(baseRepository, "refs/heads/${baseRepository.defaultBranch}")
        val baseCommit = dataService.getCommit(baseRepository, baseRef.`object`.sha)
        val fork = repositoryService.forkRepository(baseRepository)
        createOrUpdateBranchWithCommit(fork, branch, baseCommit, title, files)
        createOrUpdatePullRequest(baseRepository, fork, branch, title, description)
    }

    override fun find(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion
    ): PullRequest? {
        val baseRepository = repositoryService.getRepository({repositoryPath.fullName})
        val branch = Branch("gradle-${gradleWrapperVersion.version}-${repositoryPath.owner}")
        val query = EnhancedPullRequestService.Query(
            base = baseRepository.defaultBranch,
            head = "${userService.user}:${branch.name}",
            start = 1,
            size = 1
        )
        return pullRequestService.query(baseRepository, query).firstOrNull()
    }

    private fun createOrUpdatePullRequest(
        baseRepository: Repository,
        headRepository: Repository,
        branch: Branch,
        pullRequestTitle: String,
        pullRequestBody: String
    ): PullRequest {
        val query = EnhancedPullRequestService.Query(
            base = baseRepository.defaultBranch,
            head = "${headRepository.owner}:${branch.name}",
            state = "open",
            start = 1,
            size = 1
        )
        val latest = pullRequestService.query(baseRepository, query).firstOrNull()
        return if (latest == null) {
            pullRequestService.createPullRequest(baseRepository, PullRequest().apply {
                title = pullRequestTitle
                body = pullRequestBody
                base = PullRequestMarker().apply { label = query.base }
                head = PullRequestMarker().apply { label = query.head }
            })
        } else {
            pullRequestService.editPullRequest(baseRepository, latest.apply {
                title = pullRequestTitle
                body = pullRequestBody
            })
        }
    }

    private fun createOrUpdateBranchWithCommit(
        repository: Repository,
        branch: Branch,
        parent: Commit,
        commitMessage: String,
        files: List<GradleWrapperFile>
    ): Reference {
        val existentRef = nullIfNotFound {
            dataService.getReference(repository, branch.ref)
        }
        return if (existentRef == null) {
            dataService.createReference(repository, Reference().apply {
                ref = branch.ref
                `object` = TypedResource().apply {
                    sha = createCommit(repository, parent, commitMessage, files).sha
                }
            })
        } else {
            val existentRefParent = dataService.getCommit(repository, existentRef.`object`.sha).parents.firstOrNull()
            val needToUpdateBranch = existentRefParent?.sha == parent.sha
            if (!needToUpdateBranch) {
                dataService.editReference(repository, Reference().apply {
                    ref = branch.ref
                    `object` = TypedResource().apply {
                        sha = createCommit(repository, parent, commitMessage, files).sha
                    }
                }, true)
            } else {
                existentRef
            }
        }
    }

    private fun createCommit(
        repository: Repository,
        parent: Commit,
        commitMessage: String,
        files: List<GradleWrapperFile>
    ) =
        dataService.createCommit(repository, Commit().apply {
            message = commitMessage
            parents = listOf(parent)
            tree = dataService.createTree(repository,
                files.map { file ->
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

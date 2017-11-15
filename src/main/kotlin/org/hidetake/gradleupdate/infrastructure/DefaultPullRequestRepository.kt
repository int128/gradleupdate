package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.*
import org.eclipse.egit.github.core.service.DataService
import org.eclipse.egit.github.core.service.RepositoryService
import org.eclipse.egit.github.core.service.UserService
import org.hidetake.gradleupdate.domain.*
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedPullRequestService

@org.springframework.stereotype.Repository
class DefaultPullRequestRepository(client: SystemGitHubClient) : PullRequestRepository {
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
        val headRepository = repositoryService.forkRepository(baseRepository)
        createOrUpdateBranchWithCommit(baseRepository, headRepository, branch, title, files)
        createOrUpdatePullRequest(baseRepository, headRepository, branch, title, description)
    }

    override fun find(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion
    ): PullRequestForUpdate? {
        val baseRepository = repositoryService.getRepository({repositoryPath.fullName})
        val headBranch = Branch("gradle-${gradleWrapperVersion.version}-${repositoryPath.owner}")
        val query = EnhancedPullRequestService.Query(
            base = baseRepository.defaultBranch,
            head = "${userService.user.login}:${headBranch.name}",
            start = 1,
            size = 1
        )
        return pullRequestService.query(baseRepository, query).firstOrNull()?.let { pullRequest ->
            val parentOfHeadRef = dataService.getCommit(pullRequest.head.repo, pullRequest.head.sha).parents.firstOrNull()
            val headRefIsUpToDate = parentOfHeadRef?.sha == pullRequest.base.sha
            val state = when {
                pullRequest.state == "open" &&  headRefIsUpToDate -> PullRequestForUpdate.State.OPEN_BRANCH_UP_TO_DATE
                pullRequest.state == "open" && !headRefIsUpToDate -> PullRequestForUpdate.State.OPEN_BRANCH_OUT_OF_DATE
                pullRequest.isMerged -> PullRequestForUpdate.State.MERGED
                else -> PullRequestForUpdate.State.CLOSED
            }
            PullRequestForUpdate(state, gradleWrapperVersion, pullRequest)
        }
    }

    private fun createOrUpdatePullRequest(
        baseRepository: Repository,
        headRepository: Repository,
        headBranch: Branch,
        pullRequestTitle: String,
        pullRequestBody: String
    ): PullRequest {
        val query = EnhancedPullRequestService.Query(
            base = baseRepository.defaultBranch,
            head = "${headRepository.owner.login}:${headBranch.name}",
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
        baseRepository: Repository,
        headRepository: Repository,
        headBranch: Branch,
        commitMessage: String,
        files: List<GradleWrapperFile>
    ): Reference {
        val baseRef = dataService.getReference(baseRepository, Branch(baseRepository.defaultBranch).ref)
        val baseCommit = dataService.getCommit(baseRepository, baseRef.`object`.sha)
        val headRef = nullIfNotFound { dataService.getReference(headRepository, headBranch.ref) }
        return if (headRef == null) {
            dataService.createReference(headRepository, Reference().apply {
                ref = headBranch.ref
                `object` = TypedResource().apply {
                    sha = createCommit(headRepository, baseCommit, commitMessage, files).sha
                }
            })
        } else {
            val parentOfHeadRef = dataService.getCommit(headRepository, headRef.`object`.sha).parents.firstOrNull()
            val headRefIsUpToDate = parentOfHeadRef?.sha == baseCommit.sha
            if (!headRefIsUpToDate) {
                dataService.editReference(headRepository, Reference().apply {
                    ref = headBranch.ref
                    `object` = TypedResource().apply {
                        sha = createCommit(headRepository, baseCommit, commitMessage, files).sha
                    }
                }, true)
            } else {
                headRef
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

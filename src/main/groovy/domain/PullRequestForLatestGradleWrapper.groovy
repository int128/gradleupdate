package domain

import entity.PullRequestForLatestGradleWrapperTaskState
import groovy.util.logging.Log

import static entity.PullRequestForLatestGradleWrapperTaskState.State

@Log
class PullRequestForLatestGradleWrapper {

    private final GHRepository origin

    def PullRequestForLatestGradleWrapper(GHRepository origin) {
        this.origin = origin
    }

    boolean checkOwnership(String ownerToken) {
        final originByUser = GHSession.byToken(ownerToken).getRepository(origin.fullName)
        originByUser.checkPermissionPush()
    }

    void create() {
        def template = GradleWrapperTemplateRepository.get(origin.session)
        def templateVersion = template.gradleWrapper.version

        def originVersion = GradleVersion.getOrNull(origin.defaultBranch)
        if (originVersion == templateVersion) {
            reportProgress(State.AlreadyLatest, "Repository $origin already has the latest version $templateVersion")
        } else {
            reportProgress(State.Updating, "Updating from $originVersion to $templateVersion")

            def fork = origin.fork()
            def forkBranch = fork.defaultBranch.syncTo(origin.defaultBranch.sha)

            def gradleContents = template.gradleWrapper.fetchContents()
            def gradleCommit = forkBranch.commit("Gradle $templateVersion", gradleContents)
            def gradleBranch = fork.createOrResetBranch("gradle-$templateVersion-$origin.ownerName", gradleCommit.sha)

            def pullRequest = GHPullRequest.createOrUpdate(gradleBranch, origin.defaultBranch,
                    "Gradle $templateVersion", pullRequestBody(templateVersion))

            saveResult(pullRequest, originVersion, templateVersion)
            reportProgress(State.Done)
        }
    }

    void reportProgress(State state, String message = null) {
        log.info("$origin.fullName: $state: $message")
        new PullRequestForLatestGradleWrapperTaskState(
                fullName: origin.fullName,
                state: state,
                message: message,
                lastUpdated: new Date(),
        ).save()
    }

    void saveResult(GHPullRequest pullRequest, GradleVersion from, GradleVersion to) {
        new entity.PullRequestForLatestGradleWrapper(
                url: pullRequest.url,
                fullName: origin.fullName,
                createdAt: pullRequest.createdAt,
                fromVersion: from.string,
                toVersion: to.string,
        ).save()
    }

    private static pullRequestBody(GradleVersion version) {
        """
[Gradle $version](https://gradle.org/docs/$version/release-notes) is available now.
This pull request updates Gradle Wrapper files in the repository.

This pull request is created by [Gradle Update](https://gradleupdate.appspot.com).
""".trim()
    }

}

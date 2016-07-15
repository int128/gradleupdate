package domain

import groovy.util.logging.Log

@Log
class GradleUpdate {

    private final GHSession session

    def GradleUpdate(GHSession session = GHSession.defaultToken()) {
        this.session = session
    }

    GradleWrapperStatus getGradleWrapperStatusOrNull(String fullName, String branchName) {
        def repository = session.getRepository(fullName)
        def branch = branchName ? repository.getBranch(branchName) : repository.defaultBranch
        def version = GradleVersion.getOrNull(branch)
        version ? new GradleWrapperStatus(version) : null
    }

    @Lazy
    GradleWrapperTemplateRepository gradleWrapperTemplateRepository = {
        GradleWrapperTemplateRepository.get(session)
    }()

    PullRequestForLatestGradleWrapper pullRequestForLatestGradleWrapper(String fullName) {
        new PullRequestForLatestGradleWrapper(session.getRepository(fullName))
    }

}

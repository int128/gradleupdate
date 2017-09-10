package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.PullRequest

class GradleWrapperPullRequestStatus(
    val state: State,
    val pullRequest: PullRequest?
) {
    enum class State {
        NOT_EXIST,
        OPEN_BRANCH_UP_TO_DATE,
        OPEN_BRANCH_OUT_OF_DATE,
        CLOSED,
        MERGED
    }
}

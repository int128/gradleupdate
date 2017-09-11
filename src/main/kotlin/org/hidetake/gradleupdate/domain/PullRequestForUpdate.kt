package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.PullRequest

class PullRequestForUpdate(
    val state: State,
    val targetVersion: GradleWrapperVersion,
    val raw: PullRequest
) {
    enum class State {
        OPEN_BRANCH_UP_TO_DATE,
        OPEN_BRANCH_OUT_OF_DATE,
        CLOSED,
        MERGED
    }
}

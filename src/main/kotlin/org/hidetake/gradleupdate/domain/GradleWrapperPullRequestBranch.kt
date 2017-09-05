package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.User

class GradleWrapperPullRequestBranch(val name: String) {
    object Factory {
        fun create(repositoryName: String, version: GradleWrapperVersion) =
            GradleWrapperPullRequestBranch(
                "gradle-${version.version}-${repositoryName.split("/").first()}")
    }

    fun toLabel(forkUser: User) = "${forkUser.login}:$name"
}

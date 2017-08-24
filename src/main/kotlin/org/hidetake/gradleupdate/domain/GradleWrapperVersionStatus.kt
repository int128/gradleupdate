package org.hidetake.gradleupdate.domain

class GradleWrapperVersionStatus(
    val target: GradleWrapperVersion,
    val latest: GradleWrapperVersion
) {
    val upToDate = target.isNewerOrEqual(latest)
}

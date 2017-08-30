package org.hidetake.gradleupdate.domain

class GradleWrapperFile(
    val path: String,
    val executable: Boolean = false,
    val base64Content: String? = null
)

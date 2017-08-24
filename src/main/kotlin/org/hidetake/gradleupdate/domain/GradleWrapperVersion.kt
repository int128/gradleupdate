package org.hidetake.gradleupdate.domain

class GradleWrapperVersion(val version: String, val type: String) {
    fun isNewerOrEqual(another: GradleWrapperVersion): Boolean =
        version.split(".")
            .zip(another.version.split("."))
            .any { (left, right) ->
                left.toIntOrNull()?.let { leftNumber ->
                    right.toIntOrNull()?.let { rightNumber ->
                        leftNumber >= rightNumber
                    }
                } ?: (left >= right)
            }
}

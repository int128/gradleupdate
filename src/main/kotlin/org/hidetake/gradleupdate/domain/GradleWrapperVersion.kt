package org.hidetake.gradleupdate.domain

class GradleWrapperVersion(val version: String) {
    fun isNewerOrEqual(another: GradleWrapperVersion): Boolean {
        val leftElements = version.split(".")
        val rightElements = another.version.split(".")
        return leftElements.zip(rightElements).all { (left, right) ->
            left.toIntOrNull()?.let { leftNumber ->
                right.toIntOrNull()?.let { rightNumber ->
                    leftNumber >= rightNumber
                }
            } ?: (left >= right)
        }.let { newerOrEqualOnCommonElements ->
            when {
                newerOrEqualOnCommonElements and (leftElements.size >= rightElements.size) -> true
                else -> false
            }
        }
    }
}

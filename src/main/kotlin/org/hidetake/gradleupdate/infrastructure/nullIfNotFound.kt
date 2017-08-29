package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.client.RequestException

fun <T> nullIfNotFound(f: () -> T): T? =
    try {
        f()
    } catch (e: RequestException) {
        when (e.status) {
            404 -> null
            else -> throw e
        }
    }

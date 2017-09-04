package org.hidetake.gradleupdate.infrastructure.egit

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.stereotype.Repository
import java.security.MessageDigest

@Repository
class ResponseCacheRepository(private val memcacheService: MemcacheService) {
    fun find(uri: String): ResponseCache? =
        memcacheService.get(computeKey(uri))?.let { it as ResponseCache }

    fun save(uri: String, cache: ResponseCache) =
        memcacheService.put(computeKey(uri), cache)

    private fun computeKey(uri: String): ByteArray {
        val sha = MessageDigest.getInstance("SHA-256")
        sha.update(uri.toByteArray())
        sha.update(ResponseCache.serialVersionUID.toByte())
        return sha.digest()
    }
}

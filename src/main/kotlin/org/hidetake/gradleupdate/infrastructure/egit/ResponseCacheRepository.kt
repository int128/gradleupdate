package org.hidetake.gradleupdate.infrastructure.egit

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.stereotype.Repository
import java.security.MessageDigest

@Repository
class ResponseCacheRepository(private val memcacheService: MemcacheService) {
    fun find(uri: String, requestProperties: Map<String, List<String>>): ResponseCache? =
        memcacheService.get(computeKey(uri, requestProperties))?.let { it as ResponseCache }

    fun save(uri: String, requestProperties: Map<String, List<String>>, cache: ResponseCache) =
        memcacheService.put(computeKey(uri, requestProperties), cache)

    private fun computeKey(uri: String, requestProperties: Map<String, List<String>>): ByteArray {
        val sha = MessageDigest.getInstance("SHA-256")
        sha.update(uri.toByteArray())
        requestProperties.forEach { key, values ->
            sha.update(key.toByteArray())
            sha.update(values.joinToString().toByteArray())
        }
        sha.update(ResponseCache.serialVersionUID.toByte())
        return sha.digest()
    }
}

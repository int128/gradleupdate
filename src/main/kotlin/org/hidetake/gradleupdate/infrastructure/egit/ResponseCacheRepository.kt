package org.hidetake.gradleupdate.infrastructure.egit

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.stereotype.Repository

@Repository
class ResponseCacheRepository(private val memcacheService: MemcacheService) {
    fun find(key: ResponseCacheKey): ResponseCache? =
        memcacheService.get(key) as? ResponseCache

    fun save(key: ResponseCacheKey, cache: ResponseCache) =
        memcacheService.put(key, cache)
}

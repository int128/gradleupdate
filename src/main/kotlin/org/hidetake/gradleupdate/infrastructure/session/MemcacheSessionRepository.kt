package org.hidetake.gradleupdate.infrastructure.session

import com.google.appengine.api.memcache.Expiration
import com.google.appengine.api.memcache.MemcacheService
import org.slf4j.LoggerFactory
import org.springframework.session.SessionRepository
import org.springframework.stereotype.Component

@Component
class MemcacheSessionRepository(private val memcacheService: MemcacheService) : SessionRepository<MemcacheSession> {
    private val log = LoggerFactory.getLogger(javaClass)
    private val maxInactiveIntervalInSeconds: Int = 3600

    override fun createSession() = MemcacheSession().also { session ->
        session.maxInactiveIntervalInSeconds = maxInactiveIntervalInSeconds
        log.debug("createSession() = {}", session.id)
    }

    override fun save(session: MemcacheSession) {
        log.debug("save({}) with expiration {}", session.id, session.maxInactiveIntervalInSeconds)
        memcacheService.put(session.id, session, Expiration.byDeltaSeconds(session.maxInactiveIntervalInSeconds))
    }

    override fun getSession(id: String): MemcacheSession? =
        (memcacheService.get(id) as? MemcacheSession)?.also { session ->
            session.setLastAccessedTimeToNow()
        }.also { session ->
            log.debug("getSession({}) = {}", id, session?.id)
        }

    override fun delete(id: String) {
        log.debug("delete({})", id)
        memcacheService.delete(id)
    }
}

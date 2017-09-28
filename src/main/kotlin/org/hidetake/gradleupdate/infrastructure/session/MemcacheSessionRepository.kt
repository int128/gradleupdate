package org.hidetake.gradleupdate.infrastructure.session

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.session.SessionRepository
import org.springframework.stereotype.Component

@Component
class MemcacheSessionRepository(private val memcacheService: MemcacheService) : SessionRepository<MemcacheSession> {
    override fun createSession() = MemcacheSession().also { save(it) }

    override fun getSession(id: String) = memcacheService.get(id) as? MemcacheSession

    override fun delete(id: String) {
        memcacheService.delete(id)
    }

    override fun save(session: MemcacheSession) {
        memcacheService.put(session.id, session)
    }
}

package org.hidetake.gradleupdate.infrastructure.session

import org.springframework.session.ExpiringSession
import java.io.Serializable
import java.util.*

class MemcacheSession : ExpiringSession, Serializable {
    companion object {
        const val serialVersionUID: Long = 1
    }

    private val id: String = UUID.randomUUID().toString()
    private val creationTime: Long = System.currentTimeMillis()
    private var lastAccessedTime: Long = creationTime
    private var maxInactiveIntervalInSeconds: Int = 3600
    private val attributes: MutableMap<String, Any> = mutableMapOf()

    override fun getId() = id

    override fun getCreationTime() = creationTime

    override fun getLastAccessedTime() = lastAccessedTime
    override fun setLastAccessedTime(time: Long) {
        lastAccessedTime = time
    }
    fun setLastAccessedTimeToNow() {
        lastAccessedTime = System.currentTimeMillis()
    }

    override fun getMaxInactiveIntervalInSeconds() = maxInactiveIntervalInSeconds
    override fun setMaxInactiveIntervalInSeconds(interval: Int) {
        maxInactiveIntervalInSeconds = interval
    }

    override fun removeAttribute(key: String) {
        attributes.remove(key)
    }

    override fun getAttributeNames() = attributes.keys

    override fun <T> getAttribute(key: String): T? = attributes[key] as T?

    override fun setAttribute(key: String, value: Any) {
        attributes.put(key, value)
    }

    override fun isExpired() = false
}

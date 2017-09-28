package org.hidetake.gradleupdate.infrastructure.session

import org.springframework.session.ExpiringSession
import java.io.Serializable
import java.util.*

class MemcacheSession : ExpiringSession, Serializable {
    companion object {
        const val serialVersionUID: Long = 1
    }

    private val id = UUID.randomUUID().toString()
    private val attributes = mutableMapOf<String, Any>()
    private val creationTime = System.currentTimeMillis()
    private var lastAccessedTime = creationTime
    private var maxInactiveIntervalInSeconds = 3600

    override fun getId() = id
    override fun getAttributeNames() = attributes.keys
    override fun <T> getAttribute(attributeName: String): T = attributes[attributeName] as T
    override fun getCreationTime() = creationTime
    override fun getLastAccessedTime() = lastAccessedTime
    override fun getMaxInactiveIntervalInSeconds() = maxInactiveIntervalInSeconds

    override fun setAttribute(attributeName: String, attributeValue: Any) {
        attributes[attributeName] = attributeValue
    }

    override fun removeAttribute(attributeName: String) {
        attributes.remove(attributeName)
    }

    override fun setLastAccessedTime(time: Long) {
        lastAccessedTime = time
    }

    override fun setMaxInactiveIntervalInSeconds(interval: Int) {
        maxInactiveIntervalInSeconds = interval
    }

    override fun isExpired() = false
}

import service.GradleVersionService

final service = new GradleVersionService()

service.performIfNewStableReleaseIsAvailable { current ->
    log.info('Clear cache')
    memcache.clearCacheForUri('/stable/feed')
    memcache.clearCacheForUri('/rc/feed')

    log.info('Queue updating the Gradle template repository')
    defaultQueue.add(
            url: '/internal/update-gradle-template.groovy',
            params: [gradleVersion: current])
    defaultQueue.add(
            url: '/internal/update-gradle-of-repositories.groovy',
            params: [gradleVersion: current],
            countdownMillis: 1000 * 60)
}

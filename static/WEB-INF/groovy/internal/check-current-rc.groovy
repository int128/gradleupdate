import service.GradleVersionService

final service = new GradleVersionService()

service.performIfNewRcReleaseIsAvailable {
    log.info('Clear cache')
    memcache.clearCacheForUri('/rc/feed')
}

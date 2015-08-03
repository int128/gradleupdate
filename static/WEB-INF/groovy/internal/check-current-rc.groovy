import service.GradleVersionService

final service = new GradleVersionService()

service.performIfNewRcReleaseIsAvailable {
    memcache.clearCacheForUri('/rc/feed')
}

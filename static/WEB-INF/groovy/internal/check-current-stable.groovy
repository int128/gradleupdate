import service.GradleVersionService

final service = new GradleVersionService()

service.performIfNewStableReleaseIsAvailable { gradleVersion ->
    memcache.clearCacheForUri('/stable/feed')
    defaultQueue.add(url: '/internal/found-new-gradle/', params: [gradle_version: gradleVersion])
}

import gradle.VersionWatcher

final watcher = new VersionWatcher()

watcher.performIfNewStableReleaseIsAvailable { gradleVersion ->
    memcache.clearCacheForUri('/stable/feed')
    defaultQueue.add(url: '/internal/found-new-gradle/', params: [gradle_version: gradleVersion])
}

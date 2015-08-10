import gradle.VersionWatcher

final watcher = new VersionWatcher()

watcher.performIfNewStableReleaseIsAvailable { gradleVersion ->
    defaultQueue.add(url: '/internal/found-new-gradle/', params: [gradle_version: gradleVersion])
}

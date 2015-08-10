import gradle.VersionWatcher

final watcher = new VersionWatcher()

watcher.performIfNewStableReleaseIsAvailable { gradleVersion ->
    defaultQueue.add(
            url: '/internal/found-new-gradle/0.groovy',
            params: [gradle_version: gradleVersion]
    )
}

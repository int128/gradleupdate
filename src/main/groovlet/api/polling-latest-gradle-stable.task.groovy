import domain.GradleVersionWatcher

final watcher = new GradleVersionWatcher()

watcher.checkIfNewStableReleaseIsAvailable { gradleVersion ->
    defaultQueue.add(
            url: '/internal/api/found-new-gradle/index.task.groovy',
            params: [gradle_version: gradleVersion]
    )
}

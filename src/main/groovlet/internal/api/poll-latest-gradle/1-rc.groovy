import gradle.VersionWatcher

final watcher = new VersionWatcher()

watcher.performIfNewRcReleaseIsAvailable {
    // no task at this time
}

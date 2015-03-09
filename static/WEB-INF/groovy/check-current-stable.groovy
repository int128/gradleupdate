final service = new GradleService()

datastore.withTransaction {
    final last = CurrentGradleVersion.get('stable')?.version

    final fetched = service.fetchCurrentStableVersion().version

    if (last == fetched) {
        log.info("Current stable version is $fetched")
    } else {
        log.info("New stable version $fetched has been released")

        log.info('Clear the cache for Feed: /feed/stable')
        memcache.clearCacheForUri('/feed/stable')

        new CurrentGradleVersion(label: 'stable', version: fetched).save()
    }
}

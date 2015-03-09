final service = new GradleService()

datastore.withTransaction {
    final last = CurrentGradleVersion.get('nightly')?.version

    final fetched = service.fetchCurrentNightlyVersion().version

    if (last == fetched) {
        log.info("Current nightly version is $fetched")
    } else {
        log.info("New nightly version $fetched has been released")

        log.info('Clear the cache for Feed: /feed/all')
        memcache.clearCacheForUri('/feed/all')

        new CurrentGradleVersion(label: 'nightly', version: fetched).save()
    }
}

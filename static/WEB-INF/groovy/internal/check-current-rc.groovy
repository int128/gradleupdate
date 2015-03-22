import infrastructure.GradleRegistry
import model.CurrentGradleVersion

final registry = new GradleRegistry()

datastore.withTransaction {
    final last = CurrentGradleVersion.get('rc')?.version

    final fetched = registry.getCurrentReleaseCandidateRelease()?.version

    if (last == fetched) {
        log.info("Current rc version is $fetched")
    } else if (fetched) {
        log.info("New rc version $fetched has been released")

        log.info('Clear cache')
        memcache.clearCacheForUri('/rc/feed')

        new CurrentGradleVersion(label: 'rc', version: fetched).save()
    } else {
        log.info("Last rc version $last has been not active")

        log.info('Clear cache')
        memcache.clearCacheForUri('/rc/feed')

        CurrentGradleVersion.delete('rc')
    }
}

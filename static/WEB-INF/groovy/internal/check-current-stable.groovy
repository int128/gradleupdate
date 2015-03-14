import model.CurrentGradleVersion
import service.GradleVersionService

final service = new GradleVersionService()

datastore.withTransaction {
    final last = CurrentGradleVersion.get('stable')?.version

    final fetched = service.fetchCurrentStableVersion().version

    if (last == fetched) {
        log.info("Current stable version is $fetched")
    } else {
        log.info("New stable version $fetched has been released")

        log.info('Clear the cache for Feed: /feed/stable')
        memcache.clearCacheForUri('/feed/stable')

        log.info('Queue updating the Gradle template repository')
        defaultQueue.add(
                taskName: "update-gradle-template-$fetched",
                url: '/internal/update-gradle-template.groovy',
                params: ['gradle-version': fetched])

        new CurrentGradleVersion(label: 'stable', version: fetched).save()
    }
}

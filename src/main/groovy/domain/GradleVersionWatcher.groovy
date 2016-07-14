package domain

import entity.CurrentGradleVersion
import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import infrastructure.GradleRegistry

@Log
@GaelykBindings
class GradleVersionWatcher {

    private final registry = new GradleRegistry()

    def fetchStableVersion() {
        def cached = CurrentGradleVersion.get('stable')?.version
        if (cached) {
            cached
        } else {
            log.info("Fetching current stable release from Gradle registry")
            registry.fetchCurrentStableRelease().version
        }
    }

    def checkIfNewRcReleaseIsAvailable(Closure action) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('rc')?.version

            log.info("Fetching current RC release from Gradle registry")
            final current = registry.fetchCurrentReleaseCandidateRelease()?.version

            if (last == current) {
                log.info("RC version is still $current, do nothing")
            } else if (current) {
                log.info("Found the new RC $current, calling closure")
                action(current)
                log.info("Saving the new RC $current into the datastore")
                new CurrentGradleVersion(label: 'rc', version: current).save()
            } else {
                log.info("Last RC $last is no longer available, removing from the datastore")
                CurrentGradleVersion.delete('rc')
            }
        }
    }

    def checkIfNewStableReleaseIsAvailable(Closure action) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('stable')?.version

            log.info("Fetching current stable release from Gradle registry")
            final current = registry.fetchCurrentStableRelease().version

            if (last == current) {
                log.info("Stable version is still $current, do nothing")
            } else {
                log.info("Found the new stable $current, calling closure")
                action(current)
                log.info("Saving the new stable $current into the datastore")
                new CurrentGradleVersion(label: 'stable', version: current).save()
            }
        }
    }

}

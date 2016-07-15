package domain

import entity.CurrentGradleVersion
import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import infrastructure.GradleRegistry

@Log
@GaelykBindings
class LatestGradle {

    private final registry = new GradleRegistry()

    @Lazy
    GradleVersion version = {
        def saved = CurrentGradleVersion.get('stable')
        if (saved) {
            new GradleVersion(saved.version as String)
        } else {
            new GradleVersion(registry.fetchCurrentStableRelease().version as String)
        }
    }()

    void checkIfNewRcVersionIsAvailable(Closure action = null) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('rc')?.version
            final current = registry.fetchCurrentReleaseCandidateRelease()?.version

            if (last == current) {
                log.info("RC version is still $current, do nothing")
            } else if (current) {
                log.info("Found the new RC $current, calling closure")
                action?.call(current)
                log.info("Saving the new RC $current into the datastore")
                new CurrentGradleVersion(label: 'rc', version: current).save()
            } else {
                log.info("Last RC $last is no longer available, removing from the datastore")
                CurrentGradleVersion.delete('rc')
            }
        }
    }

    void checkIfNewStableVersionIsAvailable(Closure action = null) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('stable')?.version

            log.info("Fetching current stable release from Gradle registry")
            final current = registry.fetchCurrentStableRelease().version

            if (last == current) {
                log.info("Stable version is still $current, do nothing")
            } else {
                log.info("Found the new stable $current, calling closure")
                action?.call(current)
                log.info("Saving the new stable $current into the datastore")
                new CurrentGradleVersion(label: 'stable', version: current).save()
            }
        }
    }

}

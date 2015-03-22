package service

import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import infrastructure.GradleRegistry
import model.CurrentGradleVersion

@Log
@GaelykBindings
class GradleVersionService {

    private final registry = new GradleRegistry()

    def stableReleasesWithFixedIssues() {
        def versions = registry.getStableReleases()

        versions.take(1).each { version ->
            version.fixedIssues = registry.getIssuesFixedIn(version.version)
        }

        versions
    }

    def rcReleasesWithFixedIssues() {
        def versions = registry.getReleaseCandidateReleases()

        def rcFor = versions.find { it.rcFor }?.rcFor
        def fixedIssues = registry.getIssuesFixedIn(rcFor)

        versions.each { version ->
            if (version.version == rcFor) {
                version.fixedIssues = fixedIssues
            } else if (version.rcFor == rcFor) {
                version.fixedIssues = fixedIssues.findAll { it.fixedin == version.version }
            }
        }

        versions
    }

    def performIfNewRcReleaseIsAvailable(Closure closure) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('rc')?.version
            final current = registry.getCurrentReleaseCandidateRelease()?.version

            if (last == current) {
                log.info("Current rc release is $current")
            } else if (current) {
                log.info("New rc release ($current) is available now")
                closure(current)
                new CurrentGradleVersion(label: 'rc', version: current).save()
            } else {
                log.info("Last rc release ($last) is not active now")
                CurrentGradleVersion.delete('rc')
            }
        }
    }

    def performIfNewStableReleaseIsAvailable(Closure closure) {
        datastore.withTransaction {
            final last = CurrentGradleVersion.get('stable')?.version
            final current = registry.getCurrentStableRelease().version

            if (last == current) {
                log.info("Current stable version is $current")
            } else {
                log.info("New stable release ($current) is available now")
                closure(current)
                new CurrentGradleVersion(label: 'stable', version: current).save()
            }
        }
    }

}

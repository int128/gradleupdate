package service

import infrastructure.GradleRegistry

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

}

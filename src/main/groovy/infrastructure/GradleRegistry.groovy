package infrastructure

import groovyx.net.http.HttpURLClient

class GradleRegistry {

    private final client = new HttpURLClient(url: 'https://services.gradle.org')

    def getCurrentStableRelease() {
        client.request(path: '/versions/current').data
    }

    def getCurrentReleaseCandidateRelease() {
        client.request(path: '/versions/release-candidate').data
    }

    def getReleases() {
        client.request(path: '/versions/all').data as List
    }

    def getStableReleases() {
        getReleases().findAll { !it.snapshot && !it.rcFor }
    }

    def getReleaseCandidateReleases() {
        getReleases().findAll { !it.snapshot }
    }

    def getStableReleasesWithFixedIssues() {
        def versions = getStableReleases()
        versions.take(1).each { version ->
            version.fixedIssues = getIssuesFixedIn(version.version)
        }
        versions
    }

    def getReleaseCandidateReleasesWithFixedIssues() {
        def versions = getReleaseCandidateReleases()

        def rcFor = versions.find { it.rcFor }?.rcFor
        def fixedIssues = getIssuesFixedIn(rcFor)

        versions.each { version ->
            if (version.version == rcFor) {
                version.fixedIssues = fixedIssues
            } else if (version.rcFor == rcFor) {
                version.fixedIssues = fixedIssues.findAll { it.fixedin == version.version }
            }
        }

        versions
    }

    def getIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

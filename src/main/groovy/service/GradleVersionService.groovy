package service

import groovyx.net.http.HttpURLClient

class GradleVersionService {

    private final client = new HttpURLClient(url: 'https://services.gradle.org')

    def fetchCurrentStableVersion() {
        client.request(path: '/versions/current').data
    }

    def fetchCurrentReleaseCandidateVersion() {
        client.request(path: '/versions/release-candidate').data
    }

    def fetchCurrentNightlyVersion() {
        client.request(path: '/versions/nightly').data
    }

    def fetchAllVersions() {
        client.request(path: '/versions/all').data as List
    }

    def fetchStableVersions() {
        fetchAllVersions().findAll { !it.snapshot && !it.rcFor }
    }

    def fetchRcVersions() {
        fetchAllVersions().findAll { !it.snapshot }
    }

    def fetchStableVersionsWithFixedIssues() {
        def versions = fetchStableVersions()
        versions.take(1).each { version ->
            version.fixedIssues = fetchIssuesFixedIn(version.version)
        }
        versions
    }

    def fetchRcVersionsWithFixedIssues() {
        def versions = fetchRcVersions()

        def rcFor = versions.find { it.rcFor }?.rcFor
        def fixedIssues = fetchIssuesFixedIn(rcFor)

        versions.each { version ->
            if (version.version == rcFor) {
                version.fixedIssues = fixedIssues
            } else if (version.rcFor == rcFor) {
                version.fixedIssues = fixedIssues.findAll { it.fixedin == version.version }
            }
        }

        versions
    }

    def fetchIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

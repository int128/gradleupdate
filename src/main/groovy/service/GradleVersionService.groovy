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

    def fetchStableVersionsWithFixedIssues(int targetVersionCount) {
        def versions = fetchStableVersions()
        versions.take(targetVersionCount).each { version ->
            version.fixedIssues = fetchIssuesFixedIn(version.version)
        }
        versions
    }

    def fetchIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

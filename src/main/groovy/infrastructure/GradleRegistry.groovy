package infrastructure

import groovyx.net.http.HttpURLClient

class GradleRegistry {

    private final client = new HttpURLClient(url: 'https://services.gradle.org')

    def fetchCurrentStableRelease() {
        client.request(path: '/versions/current').data
    }

    def fetchCurrentReleaseCandidateRelease() {
        client.request(path: '/versions/release-candidate').data
    }

    def fetchReleases() {
        client.request(path: '/versions/all').data as List
    }

    def fetchStableReleases() {
        fetchReleases().findAll { !it.snapshot && !it.rcFor }
    }

    def fetchReleaseCandidateReleases() {
        fetchReleases().findAll { !it.snapshot }
    }

    def fetchIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

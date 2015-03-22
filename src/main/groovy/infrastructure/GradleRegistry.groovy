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

    def getIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

package infrastructure

import wslite.rest.RESTClient

class GradleRegistry {

    private final client = new RESTClient('https://services.gradle.org', new CacheAwareHTTPClient())

    def fetchCurrentStableRelease() {
        client.get(path: '/versions/current').json
    }

    def fetchCurrentReleaseCandidateRelease() {
        client.get(path: '/versions/release-candidate').json
    }

    def fetchReleases() {
        def releases = client.get(path: '/versions/all').json
        assert releases instanceof List
        releases
    }

    def fetchStableReleases() {
        fetchReleases().findAll { !it.snapshot && !it.rcFor }
    }

    def fetchReleaseCandidateReleases() {
        fetchReleases().findAll { !it.snapshot }
    }

    def fetchIssuesFixedIn(String version) {
        def issues = client.get(path: "/fixed-issues/$version").json
        assert issues instanceof List
        issues
    }
}

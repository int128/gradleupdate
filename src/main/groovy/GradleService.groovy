import groovyx.net.http.HttpURLClient

class GradleService {

    final HttpURLClient client

    def GradleService(String url = 'https://services.gradle.org') {
        client = new HttpURLClient(url: url)
    }

    def fetchCurrentStableVersion() {
        client.request(path: '/versions/current').data
    }

    def fetchCurrentReleaseCandidateVersion() {
        client.request(path: '/versions/release-candidate').data
    }

    def fetchCurrentNightlyVersion() {
        client.request(path: '/versions/nightly').data
    }

    List fetchAllVersions() {
        client.request(path: '/versions/all').data as List
    }

    List fetchStableVersions() {
        fetchAllVersions().findAll { !it.snapshot && !it.rcFor }
    }

    List fetchStableVersionsWithFixedIssues(int targetVersionCount) {
        def versions = fetchStableVersions()
        versions.take(targetVersionCount).each { version ->
            version.fixedIssues = fetchIssuesFixedIn(version.version)
        }
        versions
    }

    List fetchIssuesFixedIn(String version) {
        client.request(path: "/fixed-issues/$version").data as List
    }
}

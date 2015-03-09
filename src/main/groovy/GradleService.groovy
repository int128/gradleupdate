import groovy.transform.CompileDynamic
import groovy.transform.CompileStatic
import groovyx.net.http.HttpURLClient

@CompileStatic
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

    @CompileDynamic
    List fetchStableVersions() {
        fetchAllVersions().findAll { !it.snapshot && !it.rcFor }
    }

}

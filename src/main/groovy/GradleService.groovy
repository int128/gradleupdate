import groovy.transform.CompileStatic
import groovyx.net.http.HttpURLClient

@CompileStatic
class GradleService {

    final HttpURLClient client

    def GradleService(String url = 'https://services.gradle.org') {
        client = new HttpURLClient(url: url)
    }

    def fetchCurrentVersion() {
        client.request(path: '/versions/current').data
    }

}

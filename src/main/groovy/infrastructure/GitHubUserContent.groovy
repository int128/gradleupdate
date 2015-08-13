package infrastructure

import groovyx.net.http.ContentType
import groovyx.net.http.HttpURLClient
import util.HttpURLClientExtension

class GitHubUserContent implements HttpURLClientExtension {

    private final HttpURLClient client = new HttpURLClient(url: 'https://raw.githubusercontent.com')

    def fetch(String fullName, String branch, String path) {
        handleHttpResponseException(404: null) {
            def stream = client.request(
                    path: "/$fullName/$branch/$path",
                    contentType: ContentType.BINARY).data
            assert stream instanceof ByteArrayInputStream
            stream.bytes
        }
    }

}

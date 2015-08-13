package infrastructure

import com.google.appengine.api.memcache.MemcacheService
import groovy.transform.Canonical
import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import groovyx.net.http.ContentType
import groovyx.net.http.HttpURLClient
import util.HttpURLClientExtension

@Log
@GaelykBindings
class GitHubUserContent implements HttpURLClientExtension {

    private final HttpURLClient client = new HttpURLClient(url: 'https://raw.githubusercontent.com')

    def fetch(String fullName, String branch, String path) {
        fetch("/$fullName/$branch/$path")
    }

    def fetch(String fullPath) {
        assert memcache instanceof MemcacheService
        def cache = memcache.get(fullPath)
        assert cache == null || cache instanceof ContentCache

        handleHttpResponseException(404: null) {
            def response = client.request(
                    path: fullPath,
                    contentType: ContentType.BINARY,
                    headers: cache ? ['If-None-Match': cache.eTag] : [:])
            if (response.status == 304) {
                log.info("Got 304, serving GitHub content from memcache: $fullPath")
                cache.data
            } else {
                log.info("Got $response.status, serving GitHub content from response: $fullPath")
                def stream = response.data
                assert stream instanceof ByteArrayInputStream
                def data = stream.bytes
                def eTag = response.headers.etag
                if (eTag instanceof String) {
                    log.info("Updating cache for GitHub content: $fullPath")
                    memcache.put(fullPath, new ContentCache(eTag, data))
                }
                data
            }
        }
    }

    @Canonical
    static class ContentCache implements Serializable {
        static final long serialVersionUID = 1L
        String eTag
        byte[] data
    }

}

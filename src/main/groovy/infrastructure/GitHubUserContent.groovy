package infrastructure

import com.google.appengine.api.memcache.MemcacheService
import groovy.transform.Canonical
import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import wslite.rest.RESTClient

@Log
@GaelykBindings
class GitHubUserContent implements ErrorStatusHandler {

    final client = new RESTClient('https://raw.githubusercontent.com')

    def fetch(String fullName, String branch, String path) {
        fetch("/$fullName/$branch/$path")
    }

    def fetch(String fullPath) {
        assert memcache instanceof MemcacheService
        def cache = memcache.get(fullPath)
        def headers = [:]
        if (cache instanceof ContentCache) {
            headers += ['If-None-Match': cache.eTag]
        }

        handleHttpResponseException(404: null) {
            def response = client.get(path: fullPath, headers: headers)
            if (response.statusCode == 304) {
                log.info("Got 304, serving GitHub content from memcache: $fullPath")
                assert cache instanceof ContentCache
                cache.data
            } else {
                log.info("Got $response.statusCode, serving GitHub content from response: $fullPath")
                def data = response.data
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

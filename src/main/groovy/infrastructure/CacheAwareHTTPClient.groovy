package infrastructure

import com.google.appengine.api.memcache.MemcacheService
import groovyx.gaelyk.GaelykBindings
import wslite.http.HTTPClient
import wslite.http.HTTPMethod
import wslite.http.HTTPRequest
import wslite.http.HTTPResponse

import java.security.MessageDigest

@GaelykBindings
class CacheAwareHTTPClient extends HTTPClient {

    HTTPResponse execute(HTTPRequest request) {
        assert memcache instanceof MemcacheService
        if (request.method == HTTPMethod.GET) {
            def key = computeCacheKey(request)
            def cache = memcache.get(key)
            if (cache instanceof ContentCache) {
                request.headers.put('If-None-Match', cache.eTag)
            }

            def response = super.execute(request)
            if (response.statusCode == 304) {
                assert cache instanceof ContentCache
                cache.hit(response)
            } else if (response.statusCode == 200 && response.headers.containsKey('etag')) {
                memcache.put(key, new ContentCache(response))
            }
            response
        } else {
            super.execute(request)
        }
    }

    private static computeCacheKey(HTTPRequest request) {
        def digest = MessageDigest.getInstance('SHA-256')
        digest.update(request.url.toString().bytes)
        request.headers.each { k, v ->
            digest.update(k.toString().bytes)
            digest.update(v.toString().bytes)
        }
        digest.digest()
    }

    static class ContentCache implements Serializable {
        @SuppressWarnings("GroovyUnusedDeclaration")
        static final long serialVersionUID = 1L

        String eTag
        byte[] data
        String charset
        String contentEncoding
        String contentType
        int contentLength
        Map headers

        def ContentCache(HTTPResponse response) {
            eTag = response.headers.etag
            data = response.data
            charset = response.charset
            contentEncoding = response.contentEncoding
            contentType = response.contentType
            contentLength = response.contentLength
            headers = response.headers
        }

        void hit(HTTPResponse response) {
            response.data = data
            response.charset = charset
            response.contentEncoding = contentEncoding
            response.contentType = contentType
            response.contentLength = contentLength
            response.headers = headers
        }
    }

}

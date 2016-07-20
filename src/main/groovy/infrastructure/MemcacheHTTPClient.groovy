package infrastructure

import com.google.appengine.api.memcache.MemcacheService
import groovyx.gaelyk.GaelykBindings
import wslite.http.*

import java.security.MessageDigest

@GaelykBindings
class MemcacheHTTPClient extends HTTPClient {

    HTTPResponse execute(HTTPRequest request) {
        assert memcache instanceof MemcacheService
        if (request.method == HTTPMethod.GET) {
            def key = ContentCache.computeKey(request)
            def cache = memcache.get(key)
            if (cache instanceof ContentCache) {
                request.headers.put('If-None-Match', cache.eTag)
            }

            def response = executeInternal(request)
            if (response.statusCode == 304) {
                assert cache instanceof ContentCache
                cache.hit(response)
            } else if (response.statusCode == 200 && response.headers.containsKey('etag')) {
                memcache.put(key, new ContentCache(response))
            }
            response
        } else {
            executeInternal(request)
        }
    }

    private HTTPResponse executeInternal(HTTPRequest request) {
        if (app.env == 'Production') {
            super.execute(request)
        } else {
            try {
                super.execute(request)
            } catch (HTTPClientException e) {
                // openjdk raises exception on 404 but App Engine runtime does not
                if (e.cause instanceof FileNotFoundException) {
                    e.response
                } else {
                    throw e
                }
            }
        }
    }

    static class ContentCache implements Serializable {
        static final long serialVersionUID = 2L

        int statusCode
        String eTag
        byte[] data
        String charset
        String contentEncoding
        String contentType
        int contentLength
        Map headers

        def ContentCache(HTTPResponse response) {
            statusCode = response.statusCode
            eTag = response.headers.etag
            data = response.data
            charset = response.charset
            contentEncoding = response.contentEncoding
            contentType = response.contentType
            contentLength = response.contentLength
            headers = response.headers
        }

        void hit(HTTPResponse response) {
            response.statusCode = statusCode
            response.data = data
            response.charset = charset
            response.contentEncoding = contentEncoding
            response.contentType = contentType
            response.contentLength = contentLength
            response.headers = headers
        }

        static computeKey(HTTPRequest request) {
            def digest = MessageDigest.getInstance('SHA-256')
            digest.update(serialVersionUID as byte)
            digest.update(request.url.toString().bytes)
            request.headers.each { k, v ->
                digest.update(k.toString().bytes)
                digest.update(v.toString().bytes)
            }
            digest.digest()
        }
    }

}

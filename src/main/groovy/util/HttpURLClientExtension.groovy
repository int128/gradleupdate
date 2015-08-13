package util

import groovyx.net.http.ContentType
import groovyx.net.http.HttpResponseException

trait HttpURLClientExtension {

    def handleHttpResponseException(Map statusCodeMap, Closure closure) {
        try {
            closure()
        } catch (HttpResponseException e) {
            def statusCode = e.response.status
            if (statusCodeMap.containsKey(statusCode)) {
                def value = statusCodeMap[statusCode]
                log.info("Got status $statusCode from API but ignored as $value")
                value
            } else {
                throw e
            }
        }
    }

    def handle204NoContentWorkaround(Object value, Closure closure) {
        try {
            closure()
        } catch (NullPointerException e) {
            log.info("204 No Content caused NPE but ignored: $e.localizedMessage")
            value
        }
    }

    def requestJson(Map request) {
        client.request(request + [requestContentType: ContentType.JSON, body: [request.body, null]])
    }

}

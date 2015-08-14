package infrastructure

import groovyx.net.http.HttpResponseException

trait ErrorStatusHandler {

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

}

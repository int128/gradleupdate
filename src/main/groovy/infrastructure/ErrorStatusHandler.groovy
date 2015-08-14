package infrastructure

import wslite.http.HTTPClientException

trait ErrorStatusHandler {

    def handleHttpResponseException(Map statusCodeMap, Closure closure) {
        try {
            closure()
        } catch (HTTPClientException e) {
            if (e.response && statusCodeMap.containsKey(e.response.statusCode)) {
                def value = statusCodeMap[e.response.statusCode]
                log.info("Got status $e.response.statusCode from API but ignored as $value")
                value
            } else {
                throw e
            }
        }
    }

}

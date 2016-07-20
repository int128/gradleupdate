package infrastructure

import groovy.util.logging.Log
import wslite.http.HTTPClientException

@Log
class HTTPClientExceptionUtil {

    static <T> T nullIfResourceIsNotFound(Closure<T> closure) {
        try {
            closure()
        } catch (HTTPClientException e) {
            if (e.response?.statusCode == 404) {
                null
            } else {
                throw e
            }
        }
    }

}

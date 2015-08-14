package infrastructure

/**
 * A workaround trait for the bug of HttpURLClient that NPE is thrown if it got status 204.
 */
trait Status204Workaround {

    def handle204NoContentWorkaround(Object value, Closure closure) {
        try {
            closure()
        } catch (NullPointerException e) {
            log.info("204 No Content caused NPE but ignored: $e.localizedMessage")
            value
        }
    }

}

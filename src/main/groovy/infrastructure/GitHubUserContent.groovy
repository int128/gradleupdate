package infrastructure

import groovy.util.logging.Log
import groovyx.gaelyk.GaelykBindings
import wslite.rest.RESTClient

@Log
@GaelykBindings
class GitHubUserContent implements ErrorStatusHandler {

    final client = new RESTClient('https://raw.githubusercontent.com', new CacheAwareHTTPClient())

    def fetch(String fullName, String branch, String path) {
        fetch("/$fullName/$branch/$path")
    }

    def fetch(String fullPath) {
        handleHttpResponseException(404: null) {
            client.get(path: fullPath).data
        }
    }

}

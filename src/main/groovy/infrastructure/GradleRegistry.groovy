package infrastructure

import groovy.util.logging.Log
import wslite.rest.RESTClient

@Log
class GradleRegistry {

    private final client = new RESTClient('https://services.gradle.org', new MemcacheHTTPClient())

    def fetchCurrentStableRelease() {
        log.info('Fetching current stable version from Gradle registry')
        def response = client.get(path: '/versions/current')
        assert response.statusCode == 200
        response.json
    }

    def fetchCurrentReleaseCandidateRelease() {
        log.info('Fetching current RC version from Gradle registry')
        def response = client.get(path: '/versions/release-candidate')
        assert response.statusCode == 200
        response.json
    }

}

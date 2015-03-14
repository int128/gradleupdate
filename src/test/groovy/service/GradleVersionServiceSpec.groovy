package service

import service.GradleVersionService
import spock.lang.Specification

class GradleVersionServiceSpec extends Specification {

    def "fetchCurrentVersion() should return metadata of current version"() {
        given:
        def service = new GradleVersionService()

        when:
        def version = service.fetchCurrentStableVersion()

        then:
        version.current
        version.version =~ /[0-9\.]+/
    }

    def "fetchAllVersions() should return versions"() {
        given:
        def service = new GradleVersionService()

        when:
        def versions = service.fetchAllVersions()

        then:
        versions instanceof List
        versions.find { it.version == '2.3' }
        versions.find { it.version == '1.12-rc-2' }
    }

}

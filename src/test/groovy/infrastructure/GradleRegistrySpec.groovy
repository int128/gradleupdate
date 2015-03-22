package infrastructure

import spock.lang.Specification

class GradleRegistrySpec extends Specification {

    def "getCurrentStableRelease() should return metadata of current version"() {
        given:
        def service = new GradleRegistry()

        when:
        def version = service.getCurrentStableRelease()

        then:
        version.current
        version.version =~ /[0-9\.]+/
    }

    def "getReleases() should return versions"() {
        given:
        def service = new GradleRegistry()

        when:
        def versions = service.getReleases()

        then:
        versions instanceof List
        versions.find { it.version == '2.3' }
        versions.find { it.version == '1.12-rc-2' }
    }

}

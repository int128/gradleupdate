package infrastructure

import groovyx.gaelyk.spock.GaelykUnitSpec

class GradleRegistrySpec extends GaelykUnitSpec {

    def "getCurrentStableRelease() should return metadata of current version"() {
        given:
        def service = new GradleRegistry()

        when:
        def version = service.fetchCurrentStableRelease()

        then:
        version.current
        version.version =~ /[0-9\.]+/
    }

}

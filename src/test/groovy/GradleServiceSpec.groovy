import spock.lang.Specification

class GradleServiceSpec extends Specification {

    def "fetchCurrentVersion() should return metadata of current version"() {
        given:
        def service = new GradleService()

        when:
        def version = service.fetchCurrentVersion()

        then:
        version.current
        version.version =~ /[0-9\.]+/
    }

    def "fetchAllVersions() should return versions"() {
        given:
        def service = new GradleService()

        when:
        def versions = service.fetchAllVersions()

        then:
        versions instanceof List
        versions.find { it.version == '2.3' }
        versions.find { it.version == '1.12-rc-2' }
    }

}

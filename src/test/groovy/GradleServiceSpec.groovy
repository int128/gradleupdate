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

}

import com.google.appengine.tools.development.testing.LocalMemcacheServiceTestConfig
import com.google.appengine.tools.development.testing.LocalServiceTestHelper
import spock.lang.Shared
import spock.lang.Specification

class FeedServiceSpec extends Specification {

    @Shared
    def helper = new LocalServiceTestHelper(new LocalMemcacheServiceTestConfig())

    def setupSpec() {
        helper.setUp()
    }

    def cleanupSpec() {
        helper.tearDown()
    }

    def "allVersions should return versions"() {
        given:
        def service = new FeedService()

        when:
        def feed = service.allVersions
        def xml = new XmlSlurper().parseText(feed)

        then:
        xml.title.text()
    }

    def "stableVersions should return only stable versions"() {
        given:
        def service = new FeedService()

        when:
        def feed = service.stableVersions
        def xml = new XmlSlurper().parseText(feed)

        then:
        xml.title.text()
    }

}

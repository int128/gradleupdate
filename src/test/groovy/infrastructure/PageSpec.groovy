package infrastructure

import spock.lang.Specification

class PageSpec extends Specification {

    def "parseLinkHeader() should parse a link header"() {
        given:
        def linkHeader = '<https://x?y&p=4>; rel="next", <https://x?y&p=9>; rel="last", <https://x?y&p=1>; rel="first", <https://x?y&p=2>; rel="prev"'

        when:
        def rel = Page.parseLinkHeader(linkHeader)

        then:
        rel.size() == 4
        rel.next == 'https://x?y&p=4'
        rel.last == 'https://x?y&p=9'
        rel.first == 'https://x?y&p=1'
        rel.prev == 'https://x?y&p=2'
    }

    def "parseLinkHeader() should return empty object if empty string is given"() {
        when:
        def rel = Page.parseLinkHeader('')

        then:
        rel.size() == 0
    }

    def "parseLinkHeader() should return null if null is given"() {
        when:
        def rel = Page.parseLinkHeader(null)

        then:
        rel.size() == 0
    }

}

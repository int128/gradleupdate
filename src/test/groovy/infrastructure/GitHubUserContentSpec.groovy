package infrastructure

import com.google.appengine.api.memcache.MemcacheService
import groovyx.gaelyk.spock.GaelykUnitSpec

class GitHubUserContentSpec extends GaelykUnitSpec {

    def "fetch() should return a byte array of content"() {
        given:
        assert memcache instanceof MemcacheService
        def gitHubUserContent = new GitHubUserContent()

        when:
        def content = gitHubUserContent.fetch('octocat/Spoon-Knife', 'master', 'README.md')

        then:
        content instanceof byte[]
        content.size() > 0

        then:
        memcache.statistics.itemCount == 1
    }

    def "fetch(non-existent path) should return a byte array of content"() {
        given:
        def gitHubUserContent = new GitHubUserContent()

        when:
        def content = gitHubUserContent.fetch('octocat/Spoon-Knife', 'master', 'non-existent.file')

        then:
        content == null
    }

}

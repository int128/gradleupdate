package domain

import spock.lang.Specification

class GHRepositorySpec extends Specification {

    def "parse method should return version"() {
        given:
        def content = """#Sun Jun 29 00:47:42 JST 2014
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
distributionUrl=https://services.gradle.org/distributions/gradle-2.2.1-bin.zip"""

        when:
        def version = GHRepository.parseVersionFromGradleWrapperProperties(content)

        then:
        version == '2.2.1'
    }

    def "parse method should return null if content is invalid"() {
        given:
        def content = """invalid content"""

        when:
        def version = GHRepository.parseVersionFromGradleWrapperProperties(content)

        then:
        version == null
    }

    def "replaceGradleVersionString should replace version string to new"() {
        given:
        def content = """apply plugin: 'groovy'

task wrapper(type: Wrapper) {
    gradleVersion = '2.2.1'
}
"""

        when:
        def replaced = GHRepository.replaceGradleVersionString(content, '2.6')

        then:
        replaced == """apply plugin: 'groovy'

task wrapper(type: Wrapper) {
    gradleVersion = '2.6'
}
"""
    }

    def "replaceGradleVersionString should return as-is if no version syntax exists"() {
        given:
        def content = """apply plugin: 'groovy'

task wrapper(type: Wrapper) {
}
"""

        when:
        def replaced = GHRepository.replaceGradleVersionString(content, '2.6')

        then:
        replaced == content
    }

}

package infrastructure

import spock.lang.Specification

class GradleUpdateWorkerSpec extends Specification {

    def "parse method should return version"() {
        given:
        def content = """#Sun Jun 29 00:47:42 JST 2014
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
distributionUrl=https://services.gradle.org/distributions/gradle-2.2.1-bin.zip"""

        when:
        def version = GradleUpdateWorker.parseVersionFromGradleWrapperProperties(content)

        then:
        version == '2.2.1'
    }

    def "parse method should return null if content is invalid"() {
        given:
        def content = """invalid content"""

        when:
        def version = GradleUpdateWorker.parseVersionFromGradleWrapperProperties(content)

        then:
        version == null
    }

}

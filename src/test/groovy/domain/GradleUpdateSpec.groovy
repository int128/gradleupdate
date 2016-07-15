package domain

import groovyx.gaelyk.spock.GaelykUnitSpec

class GradleUpdateSpec extends GaelykUnitSpec {

    def 'should return non-null value on template project'() {
        given:
        def gradleUpdate = new GradleUpdate(GHSession.noToken())

        when:
        def status = gradleUpdate.getGradleWrapperStatusOrNull('int128/latest-gradle-wrapper', null)

        then:
        status != null
    }

    def 'should return unknown on non-Gradle project'() {
        given:
        def gradleUpdate = new GradleUpdate(GHSession.noToken())

        when:
        def status = gradleUpdate.getGradleWrapperStatusOrNull('int128/me', null)

        then:
        status == null
    }

}

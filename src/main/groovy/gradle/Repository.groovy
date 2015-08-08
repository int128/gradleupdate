package gradle

import groovyx.net.http.HttpResponseException
import infrastructure.GitHub

class Repository {

    final String fullName

    protected final GitHub gitHub

    def Repository(String fullName, GitHub gitHub) {
        this.fullName = fullName
        this.gitHub = gitHub
    }

    boolean queryIfHasGradleWrapper() {
        try {
            gitHub.getContent(fullName, 'gradle/wrapper/gradle-wrapper.properties')
            true
        } catch (HttpResponseException e) {
            if (e.statusCode == 404) {
                false
            } else {
                throw e
            }
        }
    }

    String queryGradleWrapperVersion() {
        try {
            def file = gitHub.getContent(fullName, 'gradle/wrapper/gradle-wrapper.properties')
            String base64 = file.content
            String content = new String(base64.decodeBase64())
            parseVersionFromGradleWrapperProperties(content)
        } catch (HttpResponseException e) {
            if (e.statusCode == 404) {
                null
            } else {
                throw e
            }
        }
    }

    static parseVersionFromGradleWrapperProperties(String content) {
        try {
            def matcher = content =~ /distributionUrl=.+?\/gradle-(.+?)-.+?\.zip/
            assert matcher
            assert matcher[0] instanceof List
            assert matcher[0].size() == 2
            matcher[0].get(1)
        } catch (AssertionError ignore) {
            null
        }
    }

    static updateVersionInBuildGradle(String content, String newVersion) {
        assert content
        content.replaceAll(~/(gradleVersion *= *['\"])[0-9a-z.-]+(['\"])/) {
            "${it[1]}$newVersion${it[2]}"
        }
    }

}

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

    String fetchGradleWrapperVersion() {
        try {
            def file = gitHub.fetchContent(fullName, 'gradle/wrapper/gradle-wrapper.properties')
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

    def createTreeForGradleWrapper(TemplateRepository templateRepository) {
        templateRepository.fetchGradleWrapperFiles().collect { file ->
            log.info("Creating a blob ${file.path} on $fullName")
            def blob = gitHub.createBlob(fullName, file.content).sha
            assert blob instanceof String

            log.info("Created ${file.path} as $blob on $fullName")
            [path: file.path, mode: file.mode, type: 'blob', sha: blob]
        }
    }

    def createTreeForBuildGradle(String gradleVersion) {
        final path = 'build.gradle'

        log.info("Fetching $path of $fullName")
        def content = gitHub.fetchContent(fullName, path).content
        assert content instanceof String

        def contentWithNewVersion = replaceGradleVersionString(
                new String(content.decodeBase64()), gradleVersion
        ).bytes.encodeBase64().toString()

        log.info("Creating a blob $path on $fullName")
        def blob = gitHub.createBlob(fullName, contentWithNewVersion).sha
        assert blob instanceof String

        log.info("Created $path as $blob on $fullName")
        [[path: path, mode: '100644', type: 'blob', sha: blob]]
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

    static replaceGradleVersionString(String content, String newVersion) {
        assert content
        content.replaceAll(~/(gradleVersion *= *['\"])[0-9a-z.-]+(['\"])/) {
            "${it[1]}$newVersion${it[2]}"
        }
    }

}

package gradle

import groovy.util.logging.Log
import infrastructure.GitHub

@Log
class Repository {

    final String fullName

    protected final GitHub gitHub

    def Repository(String fullName, GitHub gitHub) {
        this.fullName = fullName
        this.gitHub = gitHub
    }

    def getHtmlUrl() {
        "https://github.com/$fullName"
    }

    String fetchGradleWrapperVersion() {
        final path = 'gradle/wrapper/gradle-wrapper.properties'
        log.info("Fetching $path from repository $fullName")
        def file = gitHub.fetchContent(fullName, path)
        if (file == null) {
            log.info("Repository $fullName does not contain $path, maybe not Gradle project")
            null
        } else {
            def content = file.content
            assert content instanceof String
            parseVersionFromGradleWrapperProperties(new String(content.decodeBase64()))
        }
    }

    def createTreeForGradleWrapper(TemplateRepository templateRepository) {
        templateRepository.fetchGradleWrapperFiles().collect { file ->
            log.info("Creating a blob ${file.path} on repository $fullName")
            def blob = gitHub.createBlob(fullName, file.content).sha
            assert blob instanceof String

            log.info("Created ${file.path} as $blob on repository $fullName")
            [path: file.path, mode: file.mode, type: 'blob', sha: blob]
        }
    }

    def createTreeForBuildGradle(String gradleVersion) {
        final path = 'build.gradle'
        log.info("Fetching $path from repository $fullName")
        def file = gitHub.fetchContent(fullName, path)
        if (file == null) {
            log.info("Repository $fullName does not contain $path, no update needed")
            []
        } else {
            def content = file.content
            assert content instanceof String
            def contentWithNewVersion = replaceGradleVersionString(
                    new String(content.decodeBase64()), gradleVersion
            ).bytes.encodeBase64().toString()

            log.info("Creating a blob $path on repository $fullName")
            def blob = gitHub.createBlob(fullName, contentWithNewVersion).sha
            assert blob instanceof String

            log.info("Created $path as $blob on repository $fullName")
            [[path: path, mode: '100644', type: 'blob', sha: blob]]
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

    static replaceGradleVersionString(String content, String newVersion) {
        assert content
        content.replaceAll(~/(gradleVersion *= *['\"])[0-9a-z.-]+(['\"])/) {
            "${it[1]}$newVersion${it[2]}"
        }
    }

}

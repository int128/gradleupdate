package gradle

import groovy.util.logging.Log
import infrastructure.GitHub
import infrastructure.GitHubUserContent

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

    String fetchGradleWrapperVersion(String branch) {
        final path = 'gradle/wrapper/gradle-wrapper.properties'
        log.info("Fetching $path from repository $fullName")
        def content = new GitHubUserContent().fetch(fullName, branch, path)
        if (content == null) {
            log.info("Repository $fullName does not contain $path, maybe not Gradle project")
            null
        } else {
            assert content instanceof byte[]
            parseVersionFromGradleWrapperProperties(new String(content))
        }
    }

    GradleWrapperState checkIfGradleWrapperIsLatest(String branch) {
        def thisVersion = fetchGradleWrapperVersion(branch)
        if (thisVersion) {
            def latestVersion = new VersionWatcher().fetchStableVersion()
            if (thisVersion == latestVersion) {
                GradleWrapperState.UP_TO_DATE.for(thisVersion)
            } else {
                GradleWrapperState.OUT_OF_DATE.for(thisVersion)
            }
        } else {
            GradleWrapperState.NOT_FOUND
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

    static enum GradleWrapperState {
        UP_TO_DATE,
        OUT_OF_DATE,
        NOT_FOUND,

        def String currentVersion
        private def 'for'(String version) {
            currentVersion = version
            this
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

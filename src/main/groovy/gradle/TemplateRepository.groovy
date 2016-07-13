package gradle

import groovy.util.logging.Log

@Log
class TemplateRepository extends Repository {

    static final gradleWrapperFiles = [
            [path: 'gradlew', mode: '100755'],
            [path: 'gradlew.bat', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.properties', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.jar', mode: '100644'],
    ]

    def TemplateRepository() {
        super('int128/latest-gradle-wrapper')
    }

    def bumpVersion(String version) {
        log.info("Bump to Gradle wrapper $version on $fullName")
        removeBranch(version)
        cloneBranch('master', version)
    }

    def fetchGradleWrapperFiles() {
        gradleWrapperFiles.collect { file ->
            log.info("Fetching ${file.path} from repository $fullName")
            def content = gitHubUserContent.fetch(fullName, 'master', file.path)
            assert content instanceof byte[]
            def base64 = content.encodeBase64().toString()

            [path: file.path, mode: file.mode, content: base64]
        }
    }

}

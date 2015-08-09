package gradle

import groovy.util.logging.Log
import infrastructure.GitHub

@Log
class TemplateRepository extends Repository {

    static final gradleWrapperFiles = [
            [path: 'gradlew', mode: '100755'],
            [path: 'gradlew.bat', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.properties', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.jar', mode: '100644'],
    ]

    def TemplateRepository(GitHub gitHub) {
        super('int128/latest-gradle-wrapper', gitHub)
    }

    def bumpVersion(String version) {
        def branch = "update-gradle-template-$version"
        log.info("Recreating branch $branch on repository $fullName")
        gitHub.removeBranch(fullName, branch)
        gitHub.createBranch(fullName, branch, 'master')
    }

    def fetchGradleWrapperFiles() {
        gradleWrapperFiles.collect { file ->
            log.info("Fetching ${file.path} from repository $fullName")
            def content = gitHub.fetchContent(fullName, file.path).content
            assert content instanceof String

            [path: file.path, mode: file.mode, content: content]
        }
    }

}

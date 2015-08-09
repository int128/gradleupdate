package gradle

import groovy.util.logging.Log
import infrastructure.GitHub

@Log
class TemplateRepository extends Repository {

    static final repo = 'int128/latest-gradle-wrapper'

    static final gradleWrapperFiles = [
            [path: 'gradlew', mode: '100755'],
            [path: 'gradlew.bat', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.properties', mode: '100644'],
            [path: 'gradle/wrapper/gradle-wrapper.jar', mode: '100644'],
    ]

    def TemplateRepository(GitHub gitHub) {
        super(repo, gitHub)
    }

    def bumpTemplate(String version) {
        def branch = "update-gradle-template-$version"
        gitHub.removeBranch(repo, branch)
        gitHub.createBranch(repo, branch, 'master')
    }

    def fetchGradleWrapperFiles() {
        gradleWrapperFiles.collect { file ->
            log.info("Fetching ${file.path} of $fullName")
            def content = gitHub.fetchContent(fullName, file.path).content
            assert content instanceof String

            [path: file.path, mode: file.mode, content: content]
        }
    }

}

package gradle

import groovy.util.logging.Log
import infrastructure.GitHub

@Log
class TemplateRepository extends Repository {

    static final repo = 'int128/gradleupdate-worker'

    static final files = [
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

    List<Map> createTreeWithGradleWrapper(String targetRepo) {
        files.collect { file ->
            log.info("Fetching ${file.path} of $repo")
            def content = gitHub.getContent(repo, file.path).content
            assert content

            log.info("Creating a blob ${file.path} on $targetRepo")
            def blob = gitHub.createBlob(targetRepo, content).sha
            assert blob

            log.info("Created ${file.path} as $blob on $targetRepo")
            [path: file.path, mode: file.mode, type: 'blob', sha: blob]
        }
    }

}

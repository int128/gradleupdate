package gradle

import groovy.util.logging.Log
import infrastructure.Locator
import infrastructure.Locator.WithGitHub
import infrastructure.Locator.WithGitHubUserContent

@Log
class Repository implements WithGitHub, WithGitHubUserContent {

    final String fullName

    def Repository(String fullName) {
        this.fullName = fullName
    }

    def getHtmlUrl() {
        "https://github.com/$fullName"
    }

    def fetchGradleWrapperVersion(String branch) {
        log.info("Fetching Gradle wrapper version of repository $fullName:$branch")
        final path = 'gradle/wrapper/gradle-wrapper.properties'
        def content = gitHubUserContent.fetch(fullName, branch, path)
        if (content == null) {
            log.info("Repository $fullName does not contain $path, maybe not Gradle project")
            null
        } else {
            assert content instanceof byte[]
            parseVersionFromGradleWrapperProperties(new String(content))
        }
    }

    def checkIfGradleWrapperIsLatest(String branch) {
        log.info("Checking if repository $fullName has the latest Gradle wrapper")
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

    def createTreeForBuildGradle(String branch, String gradleVersion) {
        final path = 'build.gradle'
        log.info("Fetching $path from repository $fullName")
        def buildGradle = gitHubUserContent.fetch(fullName, branch, path)
        if (buildGradle == null) {
            log.info("Repository $fullName does not contain $path, no update needed")
            []
        } else {
            assert buildGradle instanceof byte[]
            def buildGradleWithNewVersion = replaceGradleVersionString(
                    new String(buildGradle), gradleVersion
            ).bytes.encodeBase64().toString()

            log.info("Creating a blob $path on repository $fullName")
            def blob = gitHub.createBlob(fullName, buildGradleWithNewVersion).sha
            assert blob instanceof String

            log.info("Created $path as $blob on repository $fullName")
            [[path: path, mode: '100644', type: 'blob', sha: blob]]
        }
    }

    def fork() {
        log.info("Creating a fork of repository $fullName")
        def fork = gitHub.fork(fullName)
        assert fork
        fork
    }

    def fetchPullRequests(Map filter) {
        log.info("Fetching pull requests for repository $fullName filtered by $filter")
        def pullRequests = gitHub.fetchPullRequests(filter, fullName)
        assert pullRequests instanceof List
        pullRequests
    }

    def createPullRequest(String base, String user, String branch, String title, String body) {
        log.info("Creating a pull request into $fullName:$base from $user:$branch")
        def pullRequest = gitHub.createPullRequest(fullName, base, "$user:$branch", title, body)
        assert pullRequest
        pullRequest
    }

    def remove() {
        log.info("Removing fork $fullName")
        gitHub.removeRepository(fullName)
    }

    def createBranch(String branch, String base, String commitMessage, List<Map> contents) {
        log.info("Creating a branch $branch from $base on repository $fullName")
        def baseRef = gitHub.fetchReference(fullName, base).object.sha
        assert baseRef instanceof String
        def tree = gitHub.fetchCommit(fullName, baseRef).tree.sha
        assert tree instanceof String
        def newTree = gitHub.createTree(fullName, tree, contents).sha
        assert newTree instanceof String
        def newCommit = gitHub.createCommit(fullName, [baseRef], newTree, commitMessage).sha
        assert newCommit instanceof String
        def newRef = gitHub.createReference(fullName, branch, newCommit).object.sha
        assert newRef instanceof String
        newRef
    }

    def removeBranch(String branch) {
        log.info("Removing branch $fullName:$branch")
        gitHub.removeBranch(fullName, branch)
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

    static fetchRepositories(String owner) {
        log.info("Fetching repositories of owner $owner")
        def repositories = Locator.gitHub.fetchRepositories(owner)
        assert repositories instanceof List
        repositories
    }

    static String parseVersionFromGradleWrapperProperties(String content) {
        assert content
        try {
            def m = content =~ /distributionUrl=.+?\/gradle-(.+?)-.+?\.zip/
            assert m
            def m0 = m[0]
            assert m0 instanceof List
            assert m0.size() == 2
            m0[1]
        } catch (AssertionError ignore) {
            null
        }
    }

    static String replaceGradleVersionString(String content, String newVersion) {
        assert content
        assert newVersion
        content.replaceAll(~/(gradleVersion *= *['\"])[0-9a-z.-]+(['\"])/) {
            "${it[1]}$newVersion${it[2]}"
        }
    }

}

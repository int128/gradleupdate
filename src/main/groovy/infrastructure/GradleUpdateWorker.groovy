package infrastructure

import groovyx.net.http.HttpResponseException

class GradleUpdateWorker {

    static final repo = 'int128/gradleupdate-worker'

    final gitHub = new GitHub()

    def bumpTemplate(String version) {
        def branch = "update-gradle-template-$version"
        gitHub.removeBranch(repo, branch)
        gitHub.createBranch(repo, branch, 'master')
    }

    def bumpUserRepository(String userRepo) {
        def branch = "update-gradle-of-$userRepo"
        gitHub.removeBranch(repo, branch)
        gitHub.createBranch(repo, branch, 'master')
    }

    String queryGradleWrapperVersion() {
        try {
            def file = gitHub.getContent(repo, 'gradle/wrapper/gradle-wrapper.properties')
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

}

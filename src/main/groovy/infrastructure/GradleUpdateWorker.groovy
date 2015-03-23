package infrastructure

import groovyx.net.http.HttpResponseException

class GradleUpdateWorker {

    final gitHub = new GitHub()

    def bumpTemplate(String version) {
        def branch = "update-gradle-template-$version"
        gitHub.removeBranch('gradleupdate/GradleUpdateWorker', branch)
        gitHub.createBranch('gradleupdate/GradleUpdateWorker', branch, 'master')
    }

    def bumpUserRepository(String repo) {
        def branch = "update-gradle-of-$repo"
        gitHub.removeBranch('gradleupdate/GradleUpdateWorker', branch)
        gitHub.createBranch('gradleupdate/GradleUpdateWorker', branch, 'master')
    }

    String queryGradleWrapperVersion() {
        try {
            def file = gitHub.getContent('gradleupdate/GradleUpdateWorker', 'gradle/wrapper/gradle-wrapper.properties')
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

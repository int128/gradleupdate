package service

import groovy.transform.TupleConstructor
import groovyx.net.http.HttpResponseException
import infrastructure.GitHub
import model.GitHubRepository

@TupleConstructor
class GitHubRepositoryService {

    final GitHub gitHub = new GitHub()

    def queryMetadata(String fullName) {
        assert fullName

        def repo = gitHub.getRepository(fullName)
        def gradleProject = checkIfBuildGradleExists(fullName)
        def gradleVersion = queryGradleWrapperVersion(fullName)

        def autoUpdate = null
        if (repo.permissions?.admin) {
            autoUpdate = GitHubRepository.get(fullName)?.autoUpdate ?: false
        }

        new GitHubRepositoryMetadata(
                fullName: fullName,
                admin: repo.permissions?.admin ?: false,
                gradleProject: gradleProject,
                gradleVersion: gradleVersion,
                autoUpdate: autoUpdate)
    }

    boolean saveMetadata(String fullName, boolean autoUpdate) {
        assert fullName

        def repo = gitHub.getRepository(fullName)
        if (repo.permissions?.admin) {
            final metadata = new GitHubRepository(fullName: fullName, autoUpdate: autoUpdate)
            metadata.save()
            true
        } else {
            false
        }
    }

    boolean checkIfBuildGradleExists(String fullName) {
        try {
            gitHub.getContent(fullName, 'build.gradle')
            true
        } catch (HttpResponseException e) {
            if (e.statusCode == 404) {
                false
            } else {
                throw e
            }
        }
    }

    String queryGradleWrapperVersion(String fullName) {
        try {
            def file = gitHub.getContent(fullName, 'gradle/wrapper/gradle-wrapper.properties')
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

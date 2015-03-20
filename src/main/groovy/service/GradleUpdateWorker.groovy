package service

import infrastructure.GitHub

class GradleUpdateWorker {

    final templateRepository = new GitHub()

    def bumpTemplate(String version) {
        def branch = "update-gradle-template-$version"
        templateRepository.removeBranch('gradleupdate/GradleUpdateWorker', branch)
        templateRepository.createBranch('gradleupdate/GradleUpdateWorker', branch, 'master')
    }

    def bumpUserRepository(String repo) {
        def branch = "update-gradle-of-$repo"
        templateRepository.removeBranch('gradleupdate/GradleUpdateWorker', branch)
        templateRepository.createBranch('gradleupdate/GradleUpdateWorker', branch, 'master')
    }

}

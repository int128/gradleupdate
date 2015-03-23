package infrastructure

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

}

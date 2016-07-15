package domain

class GradleWrapperTemplateRepository {

    final GHRepository repository

    private def GradleWrapperTemplateRepository(GHRepository repository) {
        this.repository = repository
    }

    static GradleWrapperTemplateRepository get(GHSession session) {
        new GradleWrapperTemplateRepository(session.getRepository('int128/latest-gradle-wrapper'))
    }

    @Lazy
    GradleWrapper gradleWrapper = {
        GradleWrapper.get(repository.defaultBranch)
    }()

    void updateAsync(String newVersion) {
        repository.removeBranch(newVersion)
        repository.defaultBranch.clone(newVersion)
    }

}

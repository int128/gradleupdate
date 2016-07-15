package domain

class GradleWrapper {

    final GHBranch branch

    private def GradleWrapper(GHBranch branch) {
        this.branch = branch
    }

    static GradleWrapper get(GHBranch branch) {
        new GradleWrapper(branch)
    }

    @Lazy
    GradleVersion version = { GradleVersion.get(branch) }()

    @Lazy
    GradleWrapperStatus status = { new GradleWrapperStatus(version) }()

    static final files = [
            'gradle/wrapper/gradle-wrapper.properties': '100644',
            'gradle/wrapper/gradle-wrapper.jar': '100644',
            'gradlew': '100755',
            'gradlew.bat': '100644'
    ]

    List<GHTreeContent> fetchContents() {
        files.collect { path, mode ->
            def content = GHContent.get(branch, path)
            new GHTreeContent(content.path, mode, content.base64encoded)
        }
    }

}

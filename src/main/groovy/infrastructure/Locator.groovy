package infrastructure

class Locator {

    @Lazy
    static GitHub gitHub = { new GitHub() }()

    static trait WithGitHub {
        GitHub getGitHub() {
            Locator.gitHub
        }
    }

    @Lazy
    static GitHubUserContent gitHubUserContent = { new GitHubUserContent() }()

    static trait WithGitHubUserContent {
        GitHubUserContent getGitHubUserContent() {
            Locator.gitHubUserContent
        }
    }

}

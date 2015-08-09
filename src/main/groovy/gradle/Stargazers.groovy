package gradle

import infrastructure.GitHub

class Stargazers {

    static final repo = 'int128/gradleupdate'

    private final gitHub

    static getHtmlUrl() {
        "https://github.com/$repo"
    }

    def Stargazers(GitHub gitHub) {
        this.gitHub = gitHub
    }

    List fetch() {
        gitHub.fetchStargazers(repo) as List
    }

}

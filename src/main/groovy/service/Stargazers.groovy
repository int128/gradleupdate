package service

import infrastructure.GitHub

class Stargazers {

    static final repo = 'int128/gradleupdate'

    private final gitHub

    def Stargazers(GitHub gitHub) {
        this.gitHub = gitHub
    }

    List fetch() {
        gitHub.getStargazers(repo) as List
    }

}

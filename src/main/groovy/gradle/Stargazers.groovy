package gradle

import infrastructure.GitHub

class Stargazers extends Repository {

    def Stargazers(GitHub gitHub) {
        super('int128/gradleupdate', gitHub)
    }

    List fetch() {
        gitHub.fetchStargazers(fullName) as List
    }

}

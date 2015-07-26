package service

import infrastructure.GitHub

class Stargazers {

    static final repo = 'int128/gradleupdate-api'

    final gitHub = new GitHub()

    List fetch() {
        gitHub.getStargazers(repo) as List
    }

}

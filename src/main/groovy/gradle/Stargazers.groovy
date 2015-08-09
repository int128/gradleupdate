package gradle

import groovy.util.logging.Log
import infrastructure.GitHub

@Log
class Stargazers extends Repository {

    def Stargazers(GitHub gitHub) {
        super('int128/gradleupdate', gitHub)
    }

    List fetch() {
        log.info("Fetching stargazers from repository $fullName")
        gitHub.fetchStargazers(fullName) as List
    }

}

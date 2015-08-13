package gradle

import groovy.util.logging.Log

@Log
class Stargazers extends Repository {

    def Stargazers() {
        super('int128/gradleupdate')
    }

    def fetch() {
        log.info("Fetching stargazers from repository $fullName")
        def stargazers = gitHub.fetchStargazers(fullName)
        assert stargazers instanceof List
        stargazers
    }

}

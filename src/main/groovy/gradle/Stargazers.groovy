package gradle

import groovy.transform.Canonical
import groovy.util.logging.Log

@Log
@Canonical
class Stargazers extends Repository {

    def Stargazers() {
        super('int128/gradleupdate')
    }

    def fetchFirst() {
        log.info("Fetching first page of stargazers from repository $fullName")
        def stargazers = gitHub.fetchStargazersOfFirstPage(fullName)
        assert stargazers, "Not found stargazers, maybe $fullName is moved"
        stargazers
    }

    def fetchNext(String next) {
        log.info("Fetching stargazers from $next")
        def stargazers = gitHub.fetchNextPage(next)
        assert stargazers, "Not found $next"
        stargazers
    }

}

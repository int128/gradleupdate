package gradle

import groovy.transform.Canonical
import groovy.util.logging.Log
import infrastructure.Locator

@Log
@Canonical
class Repositories implements Locator.WithGitHub {

    final String owner

    def Repositories(String owner) {
        this.owner = owner
    }

    def fetchFirst() {
        log.info("Fetching repositories of owner $owner")
        gitHub.fetchRepositoriesOfFirstPage(owner, sort: 'updated')
    }

    def fetchNext(String next) {
        log.info("Fetching repositories from $next")
        def repositories = gitHub.fetchNextPage(next)
        assert repositories, "Not found $next"
        repositories
    }

}

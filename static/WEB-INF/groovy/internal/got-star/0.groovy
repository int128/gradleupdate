import infrastructure.GitHub

import static util.RequestUtil.relativePath

final user = params.user
assert user instanceof String

final gitHub = new GitHub()

log.info("Fetching repositories of user $user")
final repositories = gitHub.getRepositories(user)

repositories.each { repo ->
    log.info("Queue checking the repository $repo.full_name")
    defaultQueue.add(
            url: relativePath(request, '1-check-version.groovy'),
            params: [full_name: repo.full_name])
}

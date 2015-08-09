import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fromRepo = params.from_repo
final fromBranch = params.from_branch
final intoRepo = params.into_repo
final gradleVersion = params.gradle_version
assert fromRepo instanceof String
assert fromBranch instanceof String
assert intoRepo instanceof String
assert gradleVersion instanceof String

final gitHub = new GitHub()

log.info("Removing the fork $fromRepo")
final removed = gitHub.removeRepository(fromRepo)
assert removed, "Fork $fromRepo not found, retrying"

log.info("Queue forking $intoRepo")
defaultQueue.add(
        url: relativePath(request, '3-fork.groovy'),
        params: [
                from_branch: fromBranch,
                into_repo: intoRepo,
                gradle_version: gradleVersion,
        ],
        countdownMillis: 1000)

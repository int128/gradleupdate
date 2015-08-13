import gradle.Repository

import static util.RequestUtil.relativePath

final fromRepo = params.from_repo
final fromBranch = params.from_branch
final intoRepo = params.into_repo
final intoBranch = params.into_branch
final gradleVersion = params.gradle_version
assert fromRepo instanceof String
assert fromBranch instanceof String
assert intoRepo instanceof String
assert intoBranch instanceof String
assert gradleVersion instanceof String

final removed = new Repository(fromRepo).remove()
assert removed, "Fork $fromRepo not found, retrying"

log.info("Queue recreating a fork of $intoRepo")
defaultQueue.add(
        url: relativePath(request, '3-fork.groovy'),
        params: [
                from_branch: fromBranch,
                into_repo: intoRepo,
                into_branch: intoBranch,
                gradle_version: gradleVersion,
        ],
        countdownMillis: 1000)

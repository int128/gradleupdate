import infrastructure.GitHub

import static util.RequestUtil.relativePath

final intoRepo = params.into_repo
final gradleVersion = params.gradle_version
assert intoRepo instanceof String
assert gradleVersion instanceof String

final gitHub = new GitHub()

log.info("Creating a fork of $intoRepo")
final fork = gitHub.fork(intoRepo)

final fromUser = fork.owner.login
final fromRepo = fork.full_name
final fromBranch = "gradle-$gradleVersion"
assert fromUser instanceof String
assert fromRepo instanceof String

final head = "$fromUser:$fromBranch"

log.info("Checking if any pull request exists from $head into $intoRepo")
final pullRequests = gitHub.getPullRequests(intoRepo, head: head, state: 'all')
assert pullRequests instanceof List
if (pullRequests) {
    log.info("Already sent pull requests ${pullRequests*.html_url}, skip")
    return
}

defaultQueue.add(
        url: relativePath(request, '2-remove-fork.groovy'),
        params: [
                from_repo: fromRepo,
                from_branch: fromBranch,
                into_repo: intoRepo,
                gradle_version: gradleVersion,
        ],
        countdownMillis: 1000)
log.info("Queue removing the fork: $fromRepo")

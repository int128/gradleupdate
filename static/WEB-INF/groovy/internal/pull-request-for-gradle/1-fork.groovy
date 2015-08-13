import gradle.Repository

import static util.RequestUtil.relativePath

final intoRepo = params.into_repo
final intoBranch = params.into_branch
final gradleVersion = params.gradle_version
assert intoRepo instanceof String
assert intoBranch instanceof String
assert gradleVersion instanceof String

final intoRepository = new Repository(intoRepo)

final fork = intoRepository.fork()
final fromUser = fork.owner.login
final fromRepo = fork.full_name
assert fromUser instanceof String
assert fromRepo instanceof String

final fromBranch = "gradle-$gradleVersion"
final head = "$fromUser:$fromBranch"

final pullRequests = intoRepository.fetchPullRequests(head: head, state: 'all')
if (pullRequests) {
    log.info("Already sent pull requests ${pullRequests*.html_url}, skip")
    return
}

log.info("No pull request found on repository $intoRepo, queue recreating a fork")
defaultQueue.add(
        url: relativePath(request, '2-remove-fork.groovy'),
        params: [
                from_repo: fromRepo,
                from_branch: fromBranch,
                into_repo: intoRepo,
                into_branch: intoBranch,
                gradle_version: gradleVersion,
        ],
        countdownMillis: 1000)

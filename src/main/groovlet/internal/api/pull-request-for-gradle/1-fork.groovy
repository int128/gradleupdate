import gradle.Repository

import static util.RequestUtil.relativePath

assert params.into_repo
assert params.into_branch
assert params.gradle_version

final intoRepository = new Repository(params.into_repo)

final fork = intoRepository.fork()
final fromUser = fork.owner.login
final fromRepo = fork.full_name
assert fromUser instanceof String
assert fromRepo instanceof String

final fromBranch = "gradle-$params.gradle_version"
final head = "$fromUser:$fromBranch"

final pullRequests = intoRepository.fetchPullRequests(head: head, state: 'all')
if (pullRequests) {
    log.info("Already sent pull requests ${pullRequests*.html_url}, skip")
    return
}

log.info("No pull request found on repository $params.into_repo, queue recreating a fork")
defaultQueue.add(
        url: relativePath(request, '2-remove-fork.groovy'),
        params: [
                from_repo: fromRepo,
                from_branch: fromBranch,
                into_repo: params.into_repo,
                into_branch: params.into_branch,
                gradle_version: params.gradle_version,
        ],
        countdownMillis: 1000)

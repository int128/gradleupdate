import gradle.Repository

import static util.RequestUtil.relativePath

assert params.from_branch
assert params.into_repo
assert params.into_branch
assert params.gradle_version

final fork = new Repository(params.into_repo).fork()
final fromUser = fork.owner.login
final fromRepo = fork.full_name
assert fromUser instanceof String
assert fromRepo instanceof String

log.info("Queue creating a branch $params.from_branch on $fromRepo")
defaultQueue.add(
        url: relativePath(request, '4-branch.groovy'),
        params: [
                from_user: fromUser,
                from_repo: fromRepo,
                from_branch: params.from_branch,
                into_repo: params.into_repo,
                into_branch: params.into_branch,
                gradle_version: params.gradle_version,
        ],
        countdownMillis: 1000)

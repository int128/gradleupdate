import gradle.Repository

import static util.RequestUtil.relativePath

assert params.from_repo
assert params.from_branch
assert params.into_repo
assert params.into_branch
assert params.gradle_version

final removed = new Repository(params.from_repo).remove()
assert removed, "Fork $params.from_repo not found, retrying"

log.info("Queue recreating a fork of $params.into_repo")
defaultQueue.add(
        url: relativePath(request, '3-fork.groovy'),
        params: [
                from_branch: params.from_branch,
                into_repo: params.into_repo,
                into_branch: params.into_branch,
                gradle_version: params.gradle_version,
        ],
        countdownMillis: 1000)

import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()

log.info("Creating a fork of $fullName")
final fork = gitHub.fork(fullName)
assert fork.full_name

log.info("Queue creating a branch on ${fork.full_name}")
defaultQueue.add(
        url: relativePath(request, '4-branch.groovy'),
        params: [
                full_name: fullName,
                into_branch: fork.default_branch,
                fork_name: fork.full_name,
                fork_owner: fork.owner.login,
                gradle_version: gradleVersion
        ],
        countdownMillis: 1000)

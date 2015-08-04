import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

final forkName = params.fork_name
assert forkName instanceof String

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()

log.info("Removing the fork: $forkName")
gitHub.deleteRepository(forkName)

log.info("Queue forking: $fullName")
defaultQueue.add(
        url: relativePath(request, '3-fork.groovy'),
        params: [full_name: fullName, gradle_version: gradleVersion],
        countdownMillis: 1000)

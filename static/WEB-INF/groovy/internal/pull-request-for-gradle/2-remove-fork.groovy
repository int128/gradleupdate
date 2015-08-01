import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
final forkName = params.fork_name
assert fullName instanceof String
assert forkName instanceof String

final gitHub = new GitHub()

log.info("Removing the fork: $forkName")
gitHub.deleteRepository(forkName)

log.info("Queue forking: $fullName")
defaultQueue.add(
        url: relativePath(request, '3-fork.groovy'),
        params: [full_name: fullName],
        countdownMillis: 1000)

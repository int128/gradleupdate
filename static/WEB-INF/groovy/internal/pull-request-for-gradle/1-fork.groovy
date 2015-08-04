import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()

log.info("Creating a fork of ${fullName}")
final fork = gitHub.fork(fullName)
final forkName = fork.full_name
assert forkName

final forkHead = "${fork.owner.login}:gradle-$gradleVersion"
final pullRequests = gitHub.getPullRequests(fullName, head: forkHead, state: 'all')
assert pullRequests instanceof List
if (pullRequests) {
    log.info("Already sent pull request ${pullRequests*.number} into $fullName, skip")
    return
}

log.info("Queue removing the fork: ${fork.full_name}")
defaultQueue.add(
        url: relativePath(request, '2-remove-fork.groovy'),
        params: [full_name: fullName, fork_name: forkName, gradle_version: gradleVersion],
        countdownMillis: 1000)
